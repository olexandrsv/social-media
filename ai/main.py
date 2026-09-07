from repository.repository import Repository
from service.service import Service
from server.server import Server

def main():
    repo = Repository()
    service = Service(repo)
    server = Server(service)
    server.run()

main()