Google suggests that on the 'flex' environment there should be only complementary components
to those of the 'standard' env.  
So this what I'm going to do and move ONLY the receive-notifications part.  
  
Why? My application has to do at least one autonomous request to Calendar API using
Event.Watch() method. App Engine does not allow that. It is stateless and only responds
to requests.  
I could of course make the first call by-hand with token expiration set to 'almost never',
but this is ugly and doesn't allow for any modifications later on. Say, for example, 
I'd like to receive notifications more often. I'd have to cancel the first subscription
using some saved receipt-like-thing for it (which I currenty discard) and then register new 
one. At the moment I'm receiving between 5 to 7 notifications about <i>something</i> changed
and I've got no idea how to cancel them now (I'm going to just change the endpoint name...).