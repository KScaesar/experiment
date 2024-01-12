import logging
import signal
import time
from dataclasses import dataclass
from multiprocessing import Process
from threading import Thread

_logger = logging.getLogger(__name__)


@dataclass
class Car:
    id: int
    is_close: bool = False

    def serve(self):
        while not self.is_close:
            _logger.info(f"serve car_{self.id}")
            time.sleep(3)

    def close(self):
        self.is_close = True
        _logger.info("")
        _logger.info(f"close car_{self.id}")


def graceful_shutdown(close_fn):
    def handler(sig_num, frame):
        close_fn()

    signal.signal(signal.SIGINT, handler)
    signal.signal(signal.SIGTERM, handler)


def serve_by_normal():
    car = Car(1)

    graceful_shutdown(car.close)
    car.serve()

    _logger.info("end")


def serve_by_thread():
    cars = [Car(1), Car(2)]
    thread_all = list()

    for car in cars:
        thread = Thread(target=car.serve)
        thread_all.append(thread)

    def close():
        for car in cars:
            car.close()

    graceful_shutdown(close)

    for thread in thread_all: thread.start()
    for thread in thread_all: thread.join()

    _logger.info("end")


def serve_by_process():
    cars = [Car(1), Car(2)]
    process_all = list()

    for car in cars:
        def serve_a():
            _logger.info("sub process serve")
            graceful_shutdown(car.close)
            car.serve()

        # 小心閉包問題
        # process = Process(target=serve_a)

        def serve_b(car):
            _logger.info("sub process serve")
            graceful_shutdown(car.close)
            car.serve()

        process = Process(target=serve_b, args=(car,))

        process_all.append(process)

    def close_x():
        _logger.info("main process close")
        time.sleep(5)

    graceful_shutdown(close_x)

    def close_y():
        for process in process_all:
            process.terminate()

    # main process 收到中斷訊號的時候
    # 就會傳遞 signal 到 sub process
    # 再次發送 terminate
    # 會導致 car.close 被發動兩次
    # graceful_shutdown(close_y)

    for process in process_all: process.start()
    for process in process_all: process.join()

    _logger.info("end")


if __name__ == '__main__':
    logging.basicConfig(
        level=logging.DEBUG,
        format='%(asctime)s - %(levelname)s - %(message)s',  # 设置日志格式
        datefmt='%Y-%m-%d %H:%M:%S'  # 设置日期时间格式
    )

    # serve_by_normal()
    # serve_by_thread()
    serve_by_process()
