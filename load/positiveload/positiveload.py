

from locust import HttpUser, TaskSet, task

class MyTasks(TaskSet):
    @task(1)
    def get_block(self):
        self.client.post("/block", json = {"block": 11499622, "show": True}) 
    
    @task(1)
    def get_transaction(self):
        self.client.post("/transaction", json = {"block": 11499622, "index": 12})
        
        
        
class MyWebsiteUser(HttpUser):
    tasks = [MyTasks]
    min_wait = 1000
    max_wait = 2000
    

# http://18.222.141.251:8080