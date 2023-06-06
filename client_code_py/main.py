# Example file showing a circle moving on screen
import pygame
import socket
import json

# pygame setup
pygame.init()
screen = pygame.display.set_mode((1280, 720))
clock = pygame.time.Clock()
running = True
dt = 0

SERVER_IP = "0.0.0.0"
SERVER_PORT = 3000
CURRENT_IP = "0.0.0.0"
CURRENT_PORT = 3001
RECV_BUF_SIZE = 2048

player_pos = pygame.Vector2(screen.get_width() / 2, screen.get_height() / 2)

this_player = {
    'udp_addr': CURRENT_IP + str(CURRENT_PORT),
    'id': None,
    'pos': (player_pos.x, player_pos.y)
    }

players = {}

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

sock = socket.socket(socket.AF_INET, # Internet
                     socket.SOCK_DGRAM) # UDP
sock.bind((CURRENT_IP, CURRENT_PORT))

# initiate connection
sock.sendto(bytes("Join", 'utf-8'), (CURRENT_IP, CURRENT_PORT))
data, addr = sock.recvfrom(RECV_BUF_SIZE)
this_player["id"] = str(data).strip('\x00')
print("Connection initiated with HOST")
print("Id recvd for player: ", this_player["id"])

while running:
    running = update_game(dt)
    
    # send this data 
    sock.sendto(bytes(json.dumps(this_player), 'utf-8'), (CURRENT_IP, CURRENT_PORT))
    
    # recv other data
    data, addr = sock.recvfrom(RECV_BUF_SIZE)

    print(players) 
    dt = clock.tick(60) / 1000

pygame.quit()
