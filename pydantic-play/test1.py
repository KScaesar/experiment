from pydantic import BaseModel


class User(BaseModel):
    name: str
    passowrd: str

    def change_pw(self, pw: str):
        data = {"passowrd": pw}
        return self.model_copy(
            update=data
        ).model_dump(
            exclude_unset=True
        )


def main():
    user = User(name="caesar", passowrd="123")
    print(user)
    print(user.change_pw("456"))
    pass


if __name__ == '__main__':
    main()
