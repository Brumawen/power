from gpiozero import LightSensor
from time import sleep

ldr = LightSensor(19,queue_len=1)
ldr.threshold = 0.1

while True:
    ldr.wait_for_light()
    print("light")
    ldr.wait_for_dark()
    print("dark")
