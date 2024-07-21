from prefect import flow, task
from prefect.task_runners import ConcurrentTaskRunner
import time

@task
def stop_at_floor(floor):
    print(f"elevator moving to floor {floor}")
    time.sleep(1)
    print(f"elevator stops on floor {floor}")

@flow(task_runner=ConcurrentTaskRunner())
def elevator():
    for floor in range(10, 0, -1):
        # stop_at_floor(floor) # sequence
        stop_at_floor.submit(floor) # concurrency

if __name__ == '__main__':
    elevator()