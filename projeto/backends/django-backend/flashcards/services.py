from flashcards.models import Collection, Flashcard, GenerationLog, User
from datetime import date
import json
import os

class CollectionService:
    @staticmethod
    def get_user_collections(user):
        return Collection.objects.filter(user=user)
    
    @staticmethod
    def create_collection(user, name, is_public=False):
        collection = Collection.objects.create(
            user=user,
            name=name,
            is_public=is_public
        )
        return collection
    
    @staticmethod
    def update_collection(collection, **kwargs):
        for key, value in kwargs.items():
            if value is not None:
                setattr(collection, key, value)
        collection.save()
        return collection
    
    @staticmethod
    def delete_collection(collection):
        collection.delete()

class FlashcardService:
    @staticmethod
    def get_collection_flashcards(collection):
        return Flashcard.objects.filter(collection=collection)
    
    @staticmethod
    def create_flashcard(collection, front, back, video_url=None, created_by_ia=False):
        flashcard = Flashcard.objects.create(
            collection=collection,
            front=front,
            back=back,
            video_url=video_url,
            created_by_ia=created_by_ia
        )
        return flashcard
    
    @staticmethod
    def delete_flashcard(flashcard):
        flashcard.delete()

class AIService:
    @staticmethod
    async def generate_flashcards(input_type: str, content: str):
        from openai import AsyncOpenAI
        
        client = AsyncOpenAI(api_key=os.getenv('OPENAI_API_KEY'))
        
        if input_type == "text":
            prompt = f"""Extract flashcard content from text:

Text:
{content}

Create 3-10 flashcard pairs. Respond ONLY with JSON:
[
  {{"front": "term", "back": "definition"}}
]
"""
        else:
            prompt = f"""Create 10 educational flashcards about: {content}

Respond ONLY with JSON:
[
  {{"front": "term", "back": "definition"}}
]
"""
        
        try:
            response = await client.chat.completions.create(
                model="gpt-3.5-turbo",
                messages=[
                    {"role": "system", "content": "You are a flashcard generation expert."},
                    {"role": "user", "content": prompt}
                ],
                temperature=0.7,
                max_tokens=2000
            )
            
            content = response.choices[0].message.content
            flashcards = json.loads(content)
            return flashcards
        except Exception as e:
            raise ValueError(f"Erro ao gerar flashcards: {str(e)}")

class RateLimitService:
    @staticmethod
    def check_daily_limit(user):
        today = date.today()
        log = GenerationLog.objects.filter(
            user=user,
            date=today
        ).first()
        
        plan = user.profile.plan
        daily_limit = 999 if plan == 'pro' else 6
        current_count = log.count if log else 0
        
        return current_count >= daily_limit, current_count, daily_limit
    
    @staticmethod
    def increment_generation_count(user):
        today = date.today()
        log, created = GenerationLog.objects.get_or_create(
            user=user,
            date=today
        )
        log.count += 1
        log.save()
        return log