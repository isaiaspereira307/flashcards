from pydantic import BaseModel, EmailStr, Field
from typing import Optional
from datetime import datetime
import uuid

class RegisterSchema(BaseModel):
    email: EmailStr
    password: str = Field(..., min_length=8)

class LoginSchema(BaseModel):
    email: EmailStr
    password: str

class UserSchema(BaseModel):
    id: int
    email: str
    username: str
    plan: str

class AuthResponseSchema(BaseModel):
    token: str
    refresh_token: Optional[str] = None
    user: UserSchema


class CreateCollectionSchema(BaseModel):
    name: str = Field(..., min_length=1, max_length=255)
    is_public: Optional[bool] = False

class UpdateCollectionSchema(BaseModel):
    name: Optional[str] = None
    is_public: Optional[bool] = None

class CollectionSchema(BaseModel):
    id: uuid.UUID
    name: str
    is_public: bool
    max_cards: int
    created_at: datetime


class CreateFlashcardSchema(BaseModel):
    front: str = Field(..., min_length=1)
    back: str = Field(..., min_length=1)

class FlashcardSchema(BaseModel):
    id: uuid.UUID
    front: str
    back: str
    created_by_ia: bool
    created_at: datetime


class GenerateFlashcardsSchema(BaseModel):
    input_type: str = Field(..., pattern="^(text|topic)$")
    content: str = Field(..., min_length=1)
    collection_id: uuid.UUID


class ApiResponseSchema(BaseModel):
    success: bool
    message: str
    data: Optional[dict] = None