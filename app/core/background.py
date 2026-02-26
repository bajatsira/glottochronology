import httpx
from typing import List
from app.config import settings
from app.services.calculation import run_full_analysis_async


async def perform_and_send_result(calc_id: int, lang_ids: List[int]):
    # Выполняем расчет
    results = await run_full_analysis_async(lang_ids)

    # Отправляем результат в Go
    async with httpx.AsyncClient() as client:
        try:
            url = f"{settings.callback_url}/{calc_id}/result"
            await client.put(
                url,
                json={
                    "token": settings.SECRET_TOKEN,
                    "similarityRate": results["similarityRate"],
                    "resultYearsAgo": results["resultYearsAgo"]
                }
            )
            print(f"Успех: Результат для {calc_id} отправлен в Go")
        except Exception as e:
            print(f"Ошибка при отправке в Go: {e}")
