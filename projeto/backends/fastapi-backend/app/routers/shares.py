from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session
from app.database import get_db
from app.models import User, Collection, Share
from app.schemas import ApiResponse
from app.routers.auth import get_current_user_dependency
import uuid

router = APIRouter(prefix="/shares", tags=["shares"])

@router.post("/", response_model=ApiResponse)
async def share_collection(
    collection_id: str,
    user_email: str,
    permissions: str = "read",
    db: Session = Depends(get_db),
    current_user: User = Depends(get_current_user_dependency)
):  
    collection = db.query(Collection).filter(
        Collection.id == collection_id,
        Collection.user_id == current_user.id
    ).first()
    
    if not collection:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Coleção não encontrada"
        )
        
    shared_user = db.query(User).filter(User.email == user_email).first()
    if not shared_user:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Usuário não encontrado"
        )
    
    # Criar share
    share = Share(
        collection_id=collection_id,
        user_id=shared_user.id,
        permissions=permissions,
        share_id=str(uuid.uuid4())
    )
    
    db.add(share)
    db.commit()
    db.refresh(share)
    
    return ApiResponse(
        success=True,
        message="Coleção compartilhada com sucesso",
        data={"share_id": share.share_id}
    )