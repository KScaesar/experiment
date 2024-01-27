import pytest


@pytest.fixture(scope='session')
def django_db_setup():
    pass

# https://stackoverflow.com/questions/56549271/pytest-and-django-transactional-database/56577952#56577952
def test_create_musician(django_db_blocker):
    # from play.pkg1 import fn
    # fn()
    # print("test")

    django_db_blocker.unblock()
    from play.models import create_musician
    create_musician()
