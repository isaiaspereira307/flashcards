from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session
from app.database import get_db
from app.models import User, Collection
from app.schemas import CreateCollectionRequest, UpdateCollectionRequest, CollectionResponse, ApiResponse
from app.routers.auth import get_current_user_dependency
import uuid

router = APIRouter(prefix="/collections", tags=["collections"])

@router.get("/", response_model=ApiResponse)
async def list_collections(
    db: Session = Depends(get_db),
    current_user: User = Depends(get_current_user_dependency)
):
    collections = db.query(Collection).filter(
        Collection.user_id == current_user.id
    ).all()
    
    return ApiResponse(
        success=True,
        message="Coleções recuperadas",
        data={
            "collections": [CollectionResponse.from_orm(c).dict() for c in collections]
        }
    )

@router.post("/", response_model=ApiResponse)
async def create_collection(
    request: CreateCollectionRequest,
    db: Session = Depends(get_db),
    current_user: User = Depends(get_current_user_dependency)
):    
    collection = Collection(
        user_id=current_user.id,
        name=request.name,
        is_public=request.is_public
    )
    
    db.add(collection)
    db.commit()
    db.refresh(collection)
    
    return ApiResponse(
        success=True,
        message="Coleção criada com sucesso",
        data={"collection": CollectionResponse.from_orm(collection).dict()}
    )

@router.get("/{collection_id}", response_model=ApiResponse)
async def get_collection(
    collection_id: str,
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
    
    return ApiResponse(
        success=True,
        message="Coleção recuperada",
        data={"collection": CollectionResponse.from_orm(collection).dict()}
    )

@router.put("/{collection_id}", response_model=ApiResponse)
async def update_collection(
    collection_id: str,
    request: UpdateCollectionRequest,
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
    
    if request.name:
        collection.name = request.name
    if request.is_public is not None:
        collection.is_public = request.is_public
    
    db.commit()
    db.refresh(collection)
    
    return ApiResponse(
        success=True,
        message="Coleção atualizada",
        data={"collection": CollectionResponse.from_orm(collection).dict()}
    )

@router.delete("/{collection_id}", response_model=ApiResponse)
async def delete_collection(
    collection_id: str,
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
    
    db.delete(collection)
    db.commit()
    
    return ApiResponse(
        success=True,
        message="Coleção deletada"
    )