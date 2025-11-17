from ninja import NinjaAPI
from ninja.errors import HttpError
from flashcards.models import User, Collection, Flashcard, GenerationLog
from django.contrib.auth.hashers import make_password, check_password
from flashcards.schemas import *
from flashcards.authentication import JWTAuth, create_access_token
from flashcards.services import CollectionService, FlashcardService, AIService, RateLimitService
from datetime import date

api = NinjaAPI(title="Flashcards API", version="1.0.0")
auth = JWTAuth()

# ============ AUTH ENDPOINTS ============

@api.post("/auth/register", response=ApiResponseSchema)
def register(request, data: RegisterSchema):
    """Registrar novo usuário"""
    if User.objects.filter(email=data.email).exists():
        raise HttpError(409, "Email já registrado")
    
    user = User.objects.create_user(
        email=data.email,
        username=data.email.split('@')[0],
        password=data.password,
        plan='free'
    )
    
    return {
        "success": True,
        "message": "Usuário registrado com sucesso",
        "data": {
            "user": {
                "id": str(user.id),
                "email": user.email,
                "username": user.username,
                "plan": user.plan
            }
        }
    }

@api.post("/auth/login", response=ApiResponseSchema)
def login(request, data: LoginSchema):
    """Fazer login"""
    try:
        user = User.objects.get(email=data.email)
    except User.DoesNotExist:
        raise HttpError(401, "Credenciais inválidas")
    
    if not check_password(data.password, user.password):
        raise HttpError(401, "Credenciais inválidas")
    
    token = create_access_token(user.id, user.email, user.plan)
    
    return {
        "success": True,
        "message": "Login bem-sucedido",
        "data": {
            "token": token,
            "user": {
                "id": str(user.id),
                "email": user.email,
                "username": user.username,
                "plan": user.plan
            }
        }
    }

@api.get("/auth/me", response=ApiResponseSchema, auth=auth)
def get_me(request):
    """Obter dados do usuário autenticado"""
    if not request.user:
        raise HttpError(401, "Não autorizado")
    
    return {
        "success": True,
        "message": "Dados do usuário",
        "data": {
            "user": {
                "id": str(request.user.id),
                "email": request.user.email,
                "username": request.user.username,
                "plan": request.user.plan
            }
        }
    }

# ============ COLLECTION ENDPOINTS ============

@api.get("/collections", response=ApiResponseSchema, auth=auth)
def list_collections(request):
    if not request.user:
        raise HttpError(401, "Não autorizado")
    
    collections = CollectionService.get_user_collections(request.user)
    
    return {
        "success": True,
        "message": "Coleções recuperadas",
        "data": {
            "collections": [
                {
                    "id": str(c.id),
                    "name": c.name,
                    "is_public": c.is_public,
                    "max_cards": c.max_cards,
                    "created_at": c.created_at.isoformat()
                }
                for c in collections
            ]
        }
    }

@api.post("/collections", response=ApiResponseSchema, auth=auth)
def create_collection(request, data: CreateCollectionSchema):
    if not request.user:
        raise HttpError(401, "Não autorizado")
    
    collection = CollectionService.create_collection(
        request.user,
        data.name,
        data.is_public or False
    )
    
    return {
        "success": True,
        "message": "Coleção criada",
        "data": {
            "collection": {
                "id": str(collection.id),
                "name": collection.name,
                "is_public": collection.is_public,
                "max_cards": collection.max_cards,
                "created_at": collection.created_at.isoformat()
            }
        }
    }

@api.get("/collections/{collection_id}", response=ApiResponseSchema)
def get_collection(request, collection_id: str):
    try:
        collection = Collection.objects.get(id=collection_id)
    except Collection.DoesNotExist:
        raise HttpError(404, "Coleção não encontrada")
    
    return {
        "success": True,
        "message": "Coleção encontrada",
        "data": {
            "collection": {
                "id": str(collection.id),
                "name": collection.name,
                "is_public": collection.is_public,
                "max_cards": collection.max_cards,
                "created_at": collection.created_at.isoformat()
            }
        }
    }

@api.put("/collections/{collection_id}", response=ApiResponseSchema, auth=auth)
def update_collection(request, collection_id: str, data: UpdateCollectionSchema):
    if not request.user:
        raise HttpError(401, "Não autorizado")
    
    try:
        collection = Collection.objects.get(id=collection_id, user=request.user)
    except Collection.DoesNotExist:
        raise HttpError(404, "Coleção não encontrada")
    
    collection = CollectionService.update_collection(
        collection,
        name=data.name,
        is_public=data.is_public
    )
    
    return {
        "success": True,
        "message": "Coleção atualizada",
        "data": {
            "collection": {
                "id": str(collection.id),
                "name": collection.name,
                "is_public": collection.is_public,
                "max_cards": collection.max_cards,
                "created_at": collection.created_at.isoformat()
            }
        }
    }

@api.delete("/collections/{collection_id}", response=ApiResponseSchema, auth=auth)
def delete_collection(request, collection_id: str):
    if not request.user:
        raise HttpError(401, "Não autorizado")
    
    try:
        collection = Collection.objects.get(id=collection_id, user=request.user)
    except Collection.DoesNotExist:
        raise HttpError(404, "Coleção não encontrada")
    
    CollectionService.delete_collection(collection)
    
    return {
        "success": True,
        "message": "Coleção deletada"
    }

# ============ FLASHCARD ENDPOINTS ============

@api.get("/collections/{collection_id}/flashcards", response=ApiResponseSchema)
def list_flashcards(request, collection_id: str):
    try:
        collection = Collection.objects.get(id=collection_id)
    except Collection.DoesNotExist:
        raise HttpError(404, "Coleção não encontrada")
    
    flashcards = FlashcardService.get_collection_flashcards(collection)
    
    return {
        "success": True,
        "message": "Flashcards recuperados",
        "data": {
            "flashcards": [
                {
                    "id": str(f.id),
                    "front": f.front,
                    "back": f.back,
                    "created_by_ia": f.created_by_ia,
                    "created_at": f.created_at.isoformat()
                }
                for f in flashcards
            ]
        }
    }

@api.post("/collections/{collection_id}/flashcards", response=ApiResponseSchema, auth=auth)
def create_flashcard(request, collection_id: str, data: CreateFlashcardSchema):
    if not request.user:
        raise HttpError(401, "Não autorizado")
    
    try:
        collection = Collection.objects.get(id=collection_id, user=request.user)
    except Collection.DoesNotExist:
        raise HttpError(404, "Coleção não encontrada")
    
    flashcard = FlashcardService.create_flashcard(
        collection,
        data.front,
        data.back,
        data.video_url,
        created_by_ia=False
    )
    
    return {
        "success": True,
        "message": "Flashcard criado",
        "data": {
            "flashcard": {
                "id": str(flashcard.id),
                "front": flashcard.front,
                "back": flashcard.back,
                "video_url": flashcard.video_url,
                "created_by_ia": flashcard.created_by_ia,
                "created_at": flashcard.created_at.isoformat()
            }
        }
    }

@api.delete("/flashcards/{flashcard_id}", response=ApiResponseSchema, auth=auth)
def delete_flashcard(request, flashcard_id: str):
    if not request.user:
        raise HttpError(401, "Não autorizado")
    
    try:
        flashcard = Flashcard.objects.get(id=flashcard_id, collection__user=request.user)
    except Flashcard.DoesNotExist:
        raise HttpError(404, "Flashcard não encontrado")
    
    FlashcardService.delete_flashcard(flashcard)
    
    return {
        "success": True,
        "message": "Flashcard deletado"
    }

# ============ GENERATION ENDPOINTS ============

@api.post("/flashcards/generate", response=ApiResponseSchema, auth=auth)
async def generate_flashcards(request, data: GenerateFlashcardsSchema):
    if not request.user:
        raise HttpError(401, "Não autorizado")
    
    exceeded, current, limit = RateLimitService.check_daily_limit(request.user)
    if exceeded:
        raise HttpError(429, "Limite diário de gerações atingido")
    
    try:
        flashcards = await AIService.generate_flashcards(data.input_type, data.content)
    except ValueError as e:
        raise HttpError(500, str(e))
    
    try:
        collection = Collection.objects.get(id=data.collection_id, user=request.user)
    except Collection.DoesNotExist:
        raise HttpError(404, "Coleção não encontrada")
    
    created_flashcards = []
    for card in flashcards:
        flashcard = FlashcardService.create_flashcard(
            collection,
            card['front'],
            card['back'],
            created_by_ia=True
        )
        created_flashcards.append(flashcard)
    
    RateLimitService.increment_generation_count(request.user)
    
    return {
        "success": True,
        "message": "Flashcards gerados com sucesso",
        "data": {
            "flashcards": [
                {
                    "id": str(f.id),
                    "front": f.front,
                    "back": f.back,
                    "created_by_ia": True,
                    "created_at": f.created_at.isoformat()
                }
                for f in created_flashcards
            ],
            "count": len(created_flashcards)
        }
    }

@api.get("/flashcards/generation-logs", response=ApiResponseSchema, auth=auth)
def get_generation_logs(request):
    """Obter logs de geração do usuário"""
    if not request.user:
        raise HttpError(401, "Não autorizado")
    
    today = date.today()
    log = GenerationLog.objects.filter(
        user=request.user,
        date=today
    ).first()
    
    daily_limit = 999 if request.user.plan == 'pro' else 6
    count = log.count if log else 0
    
    return {
        "success": True,
        "message": "Logs recuperados",
        "data": {
            "generated_today": count,
            "daily_limit": daily_limit,
            "remaining": daily_limit - count
        }
    }