import re
from typing import Optional

def validate_email(email: str) -> bool:
    pattern = r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'
    return re.match(pattern, email) is not None

def validate_password(password: str) -> tuple[bool, Optional[str]]:
    if len(password) < 8:
        return False, "Senha deve ter no mínimo 8 caracteres"
    
    if not re.search(r'[a-z]', password):
        return False, "Senha deve conter letras minúsculas"
    
    if not re.search(r'[A-Z]', password):
        return False, "Senha deve conter letras maiúsculas"
    
    if not re.search(r'\d', password):
        return False, "Senha deve conter números"
    
    return True, None

def validate_collection_name(name: str) -> bool:
    return 1 <= len(name) <= 255