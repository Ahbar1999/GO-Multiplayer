# Example file showing a circle moving on screen
# import pygame
import socket
import json
import time

# pygame setup
# pygame.init()
# screen = pygame.display.set_mode((1280, 720))

# clock = pygame.time.Clock()
dt = 0

HOST_IP = "127.0.0.1"
HOST_PORT = 3000
CURRENT_IP = "127.0.0.1"
CURRENT_PORT = 3001

RECV_BUF_SIZE = 2048 

# player_pos = pygame.Vector2(400 / 2, 200 / 2)
 
this_player = {
    'udp_addr': CURRENT_IP + ":" + str(CURRENT_PORT),
    'id': None,
    'pos': (0, 0)
    }

other_players = {}

'''
def update_game(dt):
    for event in pygame.event.get():
        if event.type == pygame.QUIT: 
            return False

    screen.fill("purple")

    pygame.draw.circle(screen, "red", player_pos, 40)
    
    # move our character
    keys = pygame.key.get_pressed()
    if keys[pygame.K_w]:
        player_pos.y -= 300 * dt
    if keys[pygame.K_s]:
        player_pos.y += 300 * dt
    if keys[pygame.K_a]:
        player_pos.x -= 300 * dt
    if keys[pygame.K_d]:
        player_pos.x += 300 * dt

    pygame.display.flip()

    return True
'''

if __name__ == '__main__':
    sock = socket.socket(socket.AF_INET, # Internet
                         socket.SOCK_DGRAM) # UDP
    sock.bind((CURRENT_IP, CURRENT_PORT))

    # initiate connection
    while not this_player["id"]:
        sock.sendto("Join".encode(), (HOST_IP, HOST_PORT)) 
        data, addr = sock.recvfrom(RECV_BUF_SIZE)
        if data:
            this_player["id"] = data.decode()

    print("Connection initiated with HOST")
    print("Id recvd for player: ", this_player["id"])

    running = True
    start_time = time.time()

    while True:
        print("In the loop")
        # running = update_game(dt)
        curr_time = time.time() 
        if curr_time - start_time >= 2.0:
            start_time = curr_time 
            # send this data
            send_payload = json.dumps(this_player).encode()
            print(curr_time, send_payload)
            sock.sendto(send_payload, (HOST_IP, HOST_PORT))
            
            # recv other data
            data, addr = sock.recvfrom(RECV_BUF_SIZE)

        # print(other_players) 
        # dt = clock.tick(60) / 1000

    # pygame.quit()




