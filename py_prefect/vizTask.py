from prefect import flow, task


@task(name="Print Hello")
def print_hello(name):
    msg = f"Hello {name}!"
    print(msg)
    return msg


@task(name="Print Hello Again")
def print_hello_again(name):
    msg = f"Hello {name}!"
    print(msg)
    return msg


@flow(name="Hello Flow")
def hello_world(name="world"):
    message = print_hello(name)
    message2 = print_hello_again(message)


if __name__ == "__main__":
    hello_world.visualize()
