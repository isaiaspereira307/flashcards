from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from app.config import settings
from app.database import engine, Base
from app.routers import auth, collections, flashcards, shares, subscriptions

Base.metadata.create_all(bind=engine)

app = FastAPI(
    title="Flashcards API",
    description="API for managing flashcards, collections, user authentication, sharing, and subscriptions.",
    version="1.0.0",
)
app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:3000"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.get("/health")
async def health_check():
    return {"status": "ok"}


app.include_router(auth.router)
app.include_router(collections.router)
app.include_router(flashcards.router)
app.include_router(shares.router)
app.include_router(subscriptions.router)


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(
        "main:app",
        host="0.0.0.0",
        port=8003,
        reload=settings.DEBUG,
    )
