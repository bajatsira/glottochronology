from fastapi import FastAPI
from app.api.endpoints import router

app = FastAPI(title="N+1 Glottochronology Async Service")

app.include_router(router)

@app.get("/")
async def root():
    return {"message": "Async Service is online"}
