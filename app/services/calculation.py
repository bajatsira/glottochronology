import asyncio
import random
from typing import List, Dict


async def calculate_pair_async(lang_id: int) -> Dict:
    # Имитация долгого расчета (1-3 сек)
    await asyncio.sleep(random.uniform(1, 3))
    similarity = round(random.uniform(0.6, 0.95), 4)
    print(f"Рассчитана пара для языка ID {lang_id}: {similarity * 100}%")
    return {"lang_id": lang_id, "similarity": similarity}


async def run_full_analysis_async(lang_ids: List[int]) -> Dict:
    # Запускаем задачи ПАРАЛЛЕЛЬНО (asyncio.gather)
    tasks = [calculate_pair_async(lid) for lid in lang_ids]
    results = await asyncio.gather(*tasks)

    avg_similarity = sum(r["similarity"] for r in results) / len(results)
    years = int(1000 + (1 - avg_similarity) * 10000)

    return {
        "similarityRate": avg_similarity,
        "resultYearsAgo": years
    }
