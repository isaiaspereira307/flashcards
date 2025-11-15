from openai import AsyncOpenAI
from app.config import settings
from typing import List
import json

client = AsyncOpenAI(api_key=settings.OPENAI_API_KEY)

async def generate_flashcards_with_ai(input_type: str, content: str) -> List[dict]:
    """Gerar flashcards via OpenAI"""
    
    if input_type == "text":
        prompt = f"""You are tasked with extracting flashcard content from text.

Text:
{content}

Guidelines:
1. Identify important terms and concepts
2. Extract 3-10 flashcard pairs
3. Keep front concise (1-5 words)
4. Keep back clear and educational (under 50 words)

Respond ONLY with a valid JSON array:
[
  {{"front": "term", "back": "definition"}}
]
"""
    else:  # topic
        prompt = f"""Create educational flashcards about: {content}

Guidelines:
1. Each flashcard should have clear term and definition
2. Terms should be specific (1-5 words)
3. Definitions should be clear (under 50 words)
4. Generate exactly 10 flashcards
5. Focus on most important concepts

Respond ONLY with a valid JSON array:
[
  {{"front": "term", "back": "definition"}}
]
"""
    
    try:
        response = await client.chat.completions.create(
            model="gpt-3.5-turbo",
            messages=[
                {"role": "system", "content": "You are a helpful flashcard generation assistant."},
                {"role": "user", "content": prompt}
            ],
            temperature=0.7,
            max_tokens=1000
        )
        
        # Extrair e parsear JSON
        content = response.choices[0].message.content
        flashcards = json.loads(content)
        
        return flashcards
    except json.JSONDecodeError:
        return []
    except Exception as e:
        print(f"Erro ao gerar flashcards: {e}")
        return []