## Barrier Pattern

Its purpose is simple - put up a barrier so that nobody passes until we have all the results we need, something quite common in concurrent applications. Imagine the situation where we have a micro-services application where one service needs to compose its response by merging the responses of other micro-services. This is where the Barrier pattern can help us.

Our Barrier pattern could be a service that will block its response until it has been composed with the results returned by one or more different Go-routines (or services).
Add concurrent jobs to a new barrier instance


### There are 3 ways to add jobs to a new barrier instance:

##### Add(fn functionType)

Add adds a function to our Barrier execution queue.Only use this if you don't care about fetching the response for this job later on, and only care about error.

    barrier := NewBarrier().Add(job1).Add(job2).Add(job3)
    
##### AddN(functionName string, fn functionType)

AddN adds a function to our Barrier execution queue, along with a name to the function. This can be used to fetch the corresponding result of the function later on.

    barrier := NewBarrier().AddN("job1", job1).AddN("job2", job2).AddN("job3", job3)
    
##### AddWNameReturned(fn functionType)
AddWNameReturned adds a function to our Barrier execution queue, and passes a unique name back to the user. This can be used to fetch the corresponding result of the function later on.

    barrier := NewBarrier()
    
    job1 := barrier.AddWNameReturned(job1)
    job2 := barrier.AddWNameReturned(job2)
    job3 := barrier.AddWNameReturned(job3)
    
### Execution of jobs

##### Option 1:

    results, err := barrier.Execute()
    
Execute() returns a Go or no-go, i.e. if there was an error in any of the jobs submitted, that error is returned. If all jobs passed, then all the results are returned as a map of the job name and their corresponding result.
We can just fetch the result of a function by querying the response map returned:

    //Result of Job 1 (assuming all jobs passed)
    job1Output := results["job1"]

##### Option 2:

    results := Barrier.executeAndReturnResults()
If we want more control on each of the job's result, then we can use ExecuteAndReturnResults(), which returns an array of results for us to deal with.

    for _, result := range results { 
    //process each result
     }
