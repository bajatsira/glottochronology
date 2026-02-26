from fastapi import APIRouter, HTTPException, BackgroundTasks
from app.core.background import perform_and_send_result

router = APIRouter(prefix="/api", tags=["glotto"])


@router.post("/LangCalculate")
async def start_calculation(data: dict, background_tasks: BackgroundTasks):
    calc_id = data.get("id")
    # Если фронт не прислал список ID, используем заглушку [1, 2]
    lang_ids = data.get("lang_ids", [1, 2])

    if not calc_id:
        raise HTTPException(status_code=400, detail="ID заявки обязателен")

    # Добавляем задачу в фон (методичка требует BackgroundTasks)
    background_tasks.add_task(
        perform_and_send_result,
        calc_id=calc_id,
        lang_ids=lang_ids
    )

    return {"status": "accepted", "message": "Расчет запущен в фоне"}
