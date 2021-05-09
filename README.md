# fs-challenge
FullStory tech submission. This repo hosts the backend service responsible for consuming webhook events and creating github issues.
Currently, due to limitations on my free trial account this service only creates github issues against the `MovieSearch` repo when notes are created containing the hashtag #issue.

# Getting Started
This application is deployed using Heroku and has a variety of environment variables it needs to function. 
In order to avoid leaking these values to the public they are set in the Heroku Admin console.


# Final Thoughts / Things I would do with more time.
This service is a simple example of a cool thing that could be done using the FullStory API/Webhooks. 
If I had access to an enterprise account there would be a couple features I'd implement mainly automatic issue creation when certain events reach a threshold. E.g when there are more than 500 rage clicks in the span of 4 hours.

Due to limitations there are a couple concessions I made.
1. The route does not authorize and verify the request is coming from Fullstory.
2. Multiple notes per FullStory session create comments after the initial issue is created. This forces one note per a session which is a bit restrictive but less messy than creating individual issues for each session.
3. Regex string parsing isn't the best way to determine which notes are meant to create issues. 
   If I was creating an integration like this one from the ground up I would provide a `Create Github Issue` button 
   which could bring up a modal with more useful fields/functionality that could control this behavior better. E.g Issue Title, Assignee, Projects, Labels, and etc
4. The service lacks proper test coverage. I made the decision to not write tests for the majority of the functions due to them being mostly wrappers for API calls.
    If this was a production level application I'd write tests for every method and use dependency injection to avoid polluting the code with external library structs/to allow mocking of external library functions.
5. If this was a real integration in fullstory's platform it would obviously need to be considerably more configurable. There is a lot of hardcoded values. E.g Github username, repo username, etc.
   