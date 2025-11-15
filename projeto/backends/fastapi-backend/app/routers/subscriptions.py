from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session
from app.database import get_db
from app.models import User, UserPlan, Payment
from app.schemas import ApiResponse
from app.routers.auth import get_current_user_dependency

router = APIRouter(prefix="/subscriptions", tags=["subscriptions"])

@router.post("/upgrade", response_model=ApiResponse)
async def upgrade_to_pro(
    db: Session = Depends(get_db),
    current_user: User = Depends(get_current_user_dependency)
):  
    if current_user.plan == UserPlan.PRO:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="Usuário já é PRO"
        )
    
    current_user.plan = UserPlan.PRO
    db.commit()
    db.refresh(current_user)
    
    return ApiResponse(
        success=True,
        message="Plano atualizado para PRO",
        data={"plan": current_user.plan}
    )

@router.get("/status", response_model=ApiResponse)
async def get_subscription_status(
    db: Session = Depends(get_db),
    current_user: User = Depends(get_current_user_dependency)
):    
    return ApiResponse(
        success=True,
        message="Status da assinatura",
        data={
            "plan": current_user.plan,
            "created_at": current_user.created_at
        }
    )