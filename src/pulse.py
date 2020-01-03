from gpiozero import LED
from time import sleep

led = LED(26)

led.on()
sleep(0.1)
led.off()
