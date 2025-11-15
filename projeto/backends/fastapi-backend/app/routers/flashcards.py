from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session
from app.database import get_db
from app.models import User, Collection, Flashcard, GenerationLog
from app.schemas import CreateFlashcardRequest, FlashcardResponse, GenerateFlashcardsRequest, ApiResponse
from app.routers.auth import get_current_user_dependency
from app.services.ai_service import generate_flashcards_with_ai
from datetime import date

router = APIRouter(prefix="/flashcards", tags=["flashcards"])

@router.post("/generate", response_model=ApiResponse)
async def generate_flashcards(
    collection_id: str,
    request: GenerateFlashcardsRequest,
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
    
    flashcards = await generate_flashcards_with_ai(request.input_type, request.content)
    
    cards_to_insert = [
        Flashcard(
            collection_id=collection.id,
            front=card["front"],
            back=card["back"],
            created_by_ia=True
        )
        for card in flashcards
    ]
    
    db.add_all(cards_to_insert)

    today = date.today()
    log = db.query(GenerationLog).filter(
        GenerationLog.user_id == current_user.id,
        GenerationLog.date == today
    ).first()
    
    if log:
        log.count += len(flashcards)
    else:
        log = GenerationLog(
            user_id=current_user.id,
            date=today,
            count=len(flashcards)
        )
        db.add(log)
    
    db.commit()
    
    return ApiResponse(
        success=True,
        message="Flashcards gerados com sucesso",
        data={
            "flashcards": [card for card in flashcards],
            "count": len(flashcards)
        }
    )

@router.get("/collections/{collection_id}", response_model=ApiResponse)
async def list_flashcards(
    collection_id: str,
    db: Session = Depends(get_db),
    current_user: User = Depends(get_current_user_dependency)
):
    collection = db.query(Collection).filter(
        Collection.id == collection_id
    ).first()
    
    if not collection:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Coleção não encontrada"
        )
    
    flashcards = db.query(Flashcard).filter(
        Flashcard.collection_id == collection_id
    ).all()
    
    return ApiResponse(
        success=True,
        message="Flashcards recuperados",
        data={
            "flashcards": [FlashcardResponse.from_orm(f).dict() for f in flashcards]
        }
    )

@router.post("/collections/{collection_id}", response_model=ApiResponse)
async def create_flashcard(
    collection_id: str,
    request: CreateFlashcardRequest,
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
    
    flashcard = Flashcard(
        collection_id=collection_id,
        front=request.front,
        back=request.back,
        created_by_ia=False
    )
    
    db.add(flashcard)
    db.commit()
    db.refresh(flashcard)
    
    return ApiResponse(
        success=True,
        message="Flashcard criado",
        data={"flashcard": FlashcardResponse.from_orm(flashcard).dict()}
    )

@router.delete("/{flashcard_id}", response_model=ApiResponse)
async def delete_flashcard(
    flashcard_id: str,
    db: Session = Depends(get_db),
    current_user: User = Depends(get_current_user_dependency)
):  
    flashcard = db.query(Flashcard).filter(
        Flashcard.id == flashcard_id
    ).first()
    
    if not flashcard:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Flashcard não encontrado"
        )
    
    db.delete(flashcard)
    db.commit()
    
    return ApiResponse(
        success=True,
        message="Flashcard deletado"
    )

@router.get("/generation-logs", response_model=ApiResponse)
async def get_generation_logs(
    db: Session = Depends(get_db),
    current_user: User = Depends(get_current_user_dependency)
):
    today = date.today()
    log = db.query(GenerationLog).filter(
        GenerationLog.user_id == current_user.id,
        GenerationLog.date == today
    ).first()
    
    count = log.count if log else 0
    daily_limit = 999 if current_user.plan == "pro" else 6
    
    return ApiResponse(
        success=True,
        message="Logs recuperados",
        data={
            "generated_today": count,
            "daily_limit": daily_limit,
            "remaining": daily_limit - count
        }
    )