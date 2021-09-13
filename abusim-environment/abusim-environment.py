#!/usr/bin/env python3

from signal import signal, SIGINT
from environment import Environment

env = Environment()
signal(SIGINT, lambda *_: (print(), exit(0)))

# BEGIN user defined code
def envTempChangeConv(room):
    def changeTemp(action, variables):
        del variables # just to show that the variables are there
        if action in ['increase', 'decrease']:
            temp_sensor_agent = f'temp_{room}'
            room_temp = env.get_variables(temp_sensor_agent)['temperature']
            if action == 'increase':
                env.post_input(temp_sensor_agent, f'temperature = {room_temp + 1}')
            elif action == 'decrease':
                env.post_input(temp_sensor_agent, f'temperature = {room_temp - 1}')
    return changeTemp

env.on('conv_S1', 'action', 10, envTempChangeConv('S1'))
env.on('conv_S2', 'action', 10, envTempChangeConv('S2'))
# END user defined code

env.loop()
