from prefect import flow, task


@task
def say_hello(name):
    return f"Hello {name}!"


@task
def print_result(result):
    print(type(result))
    print(result)


@flow(name="hello-flow")
def hello_world():
    future = say_hello.submit("Marvin")
    print_result.submit(future)

if __name__ == '__main__':
    hello_world()
