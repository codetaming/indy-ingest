# Ingest Design
## Design Philosophy
In the early 1970’s [Alan Kay](https://en.wikipedia.org/wiki/Alan_Kay), a computer scientist working Xerox PARC proposed a perspective of thinking about large-scale software systems. His background in Molecular Biology inspired him to consider software as a biological system as a means of enabling scaling and evolution. 

He argued that if you build an passenger aircraft and a year later require an aircraft with twice the capacity you cannot just make it bigger by adding bits on. You have to redesign it and start again and this was the same problem he was seeing with software systems. On the other hand biological systems can grow and evolve without issue.

He proposed treating software as a collection of self-contained “cells” that communicate through messages. This became the basis of Object Oriented Software design and the programming language SmallTalk, the predecessor to many modern programming languages such as Java and Swift.

With the move to cloud-based software Alan Kay’s work has come to the fore again and the idea of building highly scalable distributed systems by creating self-contained pieces and using messages between them for orchestration is a preferred approach to building cloud-based software and one I propose using on HCA DCP Ingest to maximise flexibility.

By building a system based a self-contained "cells" that perform a single task and communicate by messaging we can build a system that minimises dependencies and maximises scalability. 

## Implementation Proposal

1. Define service, high level areas of functionality, for the ingest component
2. Map out the "cells", self-contained functions, within each service and the define messages passed between them
3. Creating orchestration with dummy (no functionality) or mocks (simulated functionality) for each cell
4. Replacing dummy and mock cells with prototypes
5. Evolve prototype cells to complete implementations
6. Replace completed cells with enhanced versions if and when needed

## Advantages

* Cells can be implemented in a choice of languages and the language be changed at a cell level as limitations of the languages are reached
    * Suggested progression based on cloud provider support:
        * [Go](https://golang.org/)
        * [Node.js](https://nodejs.org/en/) (JavaScript ES6)
        * [Python 3](https://www.python.org/)
        * [Java 8](http://www.oracle.com/technetwork/java/javase/overview/java8-2100321.html) (or other JVM languages e.g. [Kotlin](https://kotlinlang.org/))
* Calls can be deployed using different plaforms and these can be changed at a cell level as a limitation of the platform is reached
    * Suggested progression:
        * Serverless ([AWS Lambda](https://aws.amazon.com/lambda), [Cloud Functions](https://cloud.google.com/functions/))
        * Container ([Kubernetes](https://kubernetes.io/), [ECS](https://aws.amazon.com/ecs))
        * Cloud VM ([EC2](https://aws.amazon.com/ec2/))
        * Local VM
* As cells are small they should take less than a day to reimplement from scatch if needed
* Orchestration between cells and services can use methods suitable for the platform:
    * Programatic Orchestration (Python Script) using a [Function Director Pattern](https://github.com/openfaas/faas/blob/master/guide/chaining_functions.md#function-director-pattern).
    * Platform Specific Orchestration (AWS Step Functions)
    * Platform Specific Messaging (AWS SNS, RabbitMQ)
    * Cross platform messaging ([Webhooks](https://en.wikipedia.org/wiki/Webhook), [Event Gateway](https://serverless.com/event-gateway/))

## References
* [The Quest to Make Code Work Like Biology Just Took A Big Step](https://www.wired.com/2016/06/chef-just-took-big-step-quest-make-code-work-like-biology/)
* [The Deep Insights of Alan Kay](http://mythz.servicestack.net/blog/2013/02/27/the-deep-insights-of-alan-kay/)
* [Seven Concurrency Models in 7 Weeks](https://www.safaribooksonline.com/library/view/seven-concurrency-models/9781941222737/0)