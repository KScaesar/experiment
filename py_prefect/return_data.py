from prefect import flow, task


@task
def add_one(x):
    return x + 1


@flow
def my_flow_1():
    result = add_one(1, return_state=True)  # return int
    print(result)
    print(result.result())

@flow
def my_flow_2():
    result = add_one.submit(1)  # return int
    print(result)
    print(result.result())
    print(result.wait())


if __name__ == '__main__':
    my_flow_2()
