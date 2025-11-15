from fastapi import APIRouter, Depends, HTTPException, status, Header
from sqlalchemy.orm import Session
from app.database import get_db
from app.models import User, UserPlan
from app.schemas import RegisterRequest, LoginRequest, UserResponse, ApiResponse
from app.security import hash_password, verify_password, create_access_token, decode_access_token, get_current_user_from_token
from typing import Optional

router = APIRouter(prefix="/auth", tags=["auth"])

def get_current_user_dependency(authorization: Optional[str] = Header(None), db: Session = Depends(get_db)):
    if not authorization:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Token não fornecido"
        )
    
    try:
        scheme, token = authorization.split()
        if scheme.lower() != "bearer":
            raise ValueError("Scheme inválido")
    except ValueError:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Formato de token inválido"
        )
    
    user_id = get_current_user_from_token(token)
    
    if not user_id:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Token inválido ou expirado"
        )
    
    user = db.query(User).filter(User.id == user_id).first()
    if not user:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Usuário não encontrado"
        )
    
    return user

@router.post("/register", response_model=ApiResponse)
async def register(request: RegisterRequest, db: Session = Depends(get_db)):
    existing_user = db.query(User).filter(User.email == request.email).first()
    if existing_user:
        raise HTTPException(
            status_code=status.HTTP_409_CONFLICT,
            detail="Email já registrado"
        )
    
    # Criar novo usuário
    user = User(
        email=request.email,
        password_hash=hash_password(request.password),
        plan=UserPlan.FREE
    )
    
    db.add(user)
    db.commit()
    db.refresh(user)
    
    return ApiResponse(
        success=True,
        message="Usuário registrado com sucesso",
        data={
            "user": UserResponse.from_orm(user).dict()
        }
    )

@router.post("/login", response_model=ApiResponse)
async def login(request: LoginRequest, db: Session = Depends(get_db)):
    user = db.query(User).filter(User.email == request.email).first()
    if not user or not verify_password(request.password, user.password_hash):
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Credenciais inválidas"
        )
    
    # Gerar token
    token = create_access_token(
        str(user.id),
        user.email,
        user.plan
    )
    
    return ApiResponse(
        success=True,
        message="Login bem-sucedido",
        data={
            "token": token,
            "user": UserResponse.from_orm(user).dict()
        }
    )

@router.get("/me", response_model=ApiResponse)
async def get_current_user(
    current_user: User = Depends(get_current_user_dependency),
):
    return ApiResponse(
        success=True,
        message="Dados do usuário",
        data={"user": UserResponse.from_orm(current_user).dict()}
    )

@router.post("/logout", response_model=ApiResponse)
async def logout():
    return ApiResponse(
        success=True,
        message="Logout realizado com sucesso"
    )

@router.post("/refresh", response_model=ApiResponse)
async def refresh_token(refresh_token: str):
    return ApiResponse(
        success=True,
        message="Token renovado"
    )