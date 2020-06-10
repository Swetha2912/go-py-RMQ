import time
import collections


class FPS:

    def __init__(self, name=None, start=False, preserve_history=False, print_counter=100, reset_on_print=True):
        self.name = name
        self.counter = 0
        self.print_counter = print_counter
        self.reset_on_print = reset_on_print

        if start:
            self.start_time = time.time()
        else:
            self.start_time = None

        if preserve_history:
            self.history = collections.deque(maxlen=10)
        else:
            self.history = None

    def start(self):
        self.start_time = time.time()

    def reset(self):
        self.start_time = time.time()
        self.counter = 0

    def plus(self):
        self.counter += 1
        if self.print_counter is not None:
            if self.counter % self.print_counter == 0:
                self.calculate(True)

    def calculate(self, print_fps=False):
        end_time = time.time()
        elapsed_time = end_time - self.start_time
        calculated_fps = self.counter / elapsed_time

        if print_fps:
            print(self.name, "fps", calculated_fps)

        if self.history is not None:
            self.history.append(calculated_fps)

        if self.reset_on_print:
            self.counter = 0
            self.start_time = time.time()

        return calculated_fps
