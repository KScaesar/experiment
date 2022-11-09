from dataclasses import dataclass


@dataclass
class TestNode:
    value: int
    previous: 'TestNode' = None
    next: 'TestNode' = None
