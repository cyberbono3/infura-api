
from locust import HttpUser, TaskSet, task

class MyTasks(TaskSet):
    @task(1)
    def missing_showTransFlag(self):
      with  self.client.post("/block", json = {"block": 11499622}, catch_response=True) as response:
        if response.status_code == 200:
            response.success()         
    
        
    @task(1)
    def missing_block_number(self):
        with self.client.post("/block", json = {"show": True}, catch_response=True) as response:
             if response.status_code == 422:
                response.success()         
    
    @task(1)
    def negative_block_number(self):
        with self.client.post("/block", json = {"block": -11499622, "show": True}, catch_response=True) as response:
             if response.status_code == 422:
                response.success()    
                
        
    @task(1)
    def invalid_json_input_type(self):
        with self.client.post("/block", json = {"block": "11499622", "show": "True"}, catch_response=True) as response:
             if response.status_code == 400:
                response.success()    
    
               
class MyWebsiteUser(HttpUser):
    tasks = [MyTasks]
    min_wait = 1000
    max_wait = 2000
    
