import jwt
from datetime import datetime, timedelta
from django.conf import settings
from flashcards.models import User
from ninja.security import HttpBearer
from typing import Optional

class JWTAuth(HttpBearer):
    def authenticate(self, request, token):
        try:
            payload = jwt.decode(
                token,
                settings.JWT_SECRET,
                algorithms=[settings.JWT_ALGORITHM]
            )
            user_id = payload.get('sub')
            user = User.objects.get(id=user_id)
            request.user = user
            return user
        except (jwt.InvalidTokenError, User.DoesNotExist):
            return None

def create_access_token(user_id: int, email: str, plan: str) -> str:
    payload = {
        'sub': user_id,
        'email': email,
        'plan': plan,
        'exp': datetime.utcnow() + timedelta(hours=settings.JWT_EXPIRATION_HOURS),
        'iat': datetime.utcnow()
    }
    
    token = jwt.encode(
        payload,
        settings.JWT_SECRET,
        algorithm=settings.JWT_ALGORITHM
    )
    
    return token

def decode_access_token(token: str) -> Optional[dict]:
    try:
        payload = jwt.decode(
            token,
            settings.JWT_SECRET,
            algorithms=[settings.JWT_ALGORITHM]
        )
        return payload
    except jwt.ExpiredSignatureError:
        return None
    except jwt.InvalidTokenError:
        return None