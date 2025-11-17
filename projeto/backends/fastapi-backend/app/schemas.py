from pydantic import BaseModel, EmailStr, Field, field_validator
from typing import Optional
from datetime import datetime
import uuid

class RegisterRequest(BaseModel):
    email: EmailStr
    password: str = Field(..., min_length=8)

class LoginRequest(BaseModel):
    email: EmailStr
    password: str

class UserResponse(BaseModel):
    id: uuid.UUID
    email: str
    plan: str
    created_at: datetime
    
    class Config:
        from_attributes = True


class CreateCollectionRequest(BaseModel):
    name: str = Field(..., min_length=1, max_length=255)
    is_public: Optional[bool] = False

class UpdateCollectionRequest(BaseModel):
    name: Optional[str] = None
    is_public: Optional[bool] = None

class CollectionResponse(BaseModel):
    id: uuid.UUID
    name: str
    is_public: bool
    max_cards: int
    created_at: datetime
    
    class Config:
        from_attributes = True


class CreateFlashcardRequest(BaseModel):
    front: str = Field(..., min_length=1)
    back: str = Field(..., min_length=1)
    video_url: Optional[str] = Field(None, max_length=1000)

    @field_validator('video_url')
    @classmethod
    def validate_video_url(cls, v):
        if v and not ("youtube.com" in v or "vimeo.com" in v or "youtu.be" in v):
            raise ValueError("video_url must be a valid URL")
        return v

class FlashcardResponse(BaseModel):
    id: uuid.UUID
    collection_id: uuid.UUID
    front: str
    back: str
    video_url: Optional[str] = None
    created_by_ia: bool
    created_at: datetime
    
    class Config:
        from_attributes = True

class GenerateFlashcardsRequest(BaseModel):
    input_type: str = Field(..., pattern="^(text|topic)$")
    content: str = Field(..., min_length=10)


class ApiResponse(BaseModel):
    success: bool
    message: str
    data: Optional[dict] = None