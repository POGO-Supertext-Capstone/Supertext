# RFC: Client-Side Frontend Conventions

| Field | Value |
| --- | --- |
| Author | Tom Sharp |
| Date | 2026-10-25 |
| Status | Proposed |
| Decision Date | N/A |
| Deciders | N/A |
| Superseded By | N/A |

## Problem and Context

The Supertext program needs to run off a Raspberry Pi in an area that has no internet, and be able to send text to possibly hundreds of cell phones simultaneously and with low latency. Because of this, we need a frontend framework that is lightweight and efficient.

The decision needed should be able to answer the questions, how do we get text from the database to the clients, and how do we trigger slide changes, all while maintaining the lightweight and efficiency requirements.

## Options

### Option A: HTMX

HTMX is a lightweight JavaScript library that extends HTML with custom keywords. It allows for embedding endpoints directly into HTML tags without the need for custom JS or forms.

This lightweight nature would make it ideal for running with a Raspberry Pi, however there are some downsides to this approach.

**Server Reliance:** UI updates rely on the server. Take the following story: a showrunner want to trigger another slide, then each client will have to [have the server render the HTML](https://www.syncfusion.com/blogs/post/htmx-with-react-nextjs-server-driven-ui) and send it back. With possibly hundreds of people, this can get very costly very fast.

**Low Flexability:** This also connects with the server reliance, but HTMX would make it very hard to use lightweight server sent events to trigger slide changes. We couldn't just tell the client which slide to change to, we'd need to render the text and send that. That just isn't scalable.

For low participants however, this could very well be the right choice. It keeps the [frontend exceptionally small](https://dev.to/leapcell/htmx-a-simpler-frontend-revolution-forget-react-171) and is generally easier to code and create for developers.

### Option B: Pure JS / HTML

We could use pure JavaScript to handle interactions. This would be beneficial as it would give us the flexibility on how data is handled, and allow the client to take on more of the work from the server to keep latency low.

This would allow us to implement an approach where we send the entire slideshow to the client once, then use server sent events to trigger slide changes. The client could then update the DOM itself, saving the server resources.

The big downside for this approach is the encoding / decoding of the slideshows. We could do a lot to make it plug-and-play, but if some major storage change were to happen, we might need to also alter the frontend code to be able to decode it.

If we had a very large number of participants, this could be the right choice as its flexibility alone would allow us to optimize as much as we needed to get the features working without overwhelming the pi.

### Option C: Hybrid

A hybrid approach would allow us the benefits of HTMX while also allowing us the flexibility of pure JavaScript. We can use HTMX to retrieve the slideshow from the server, then use JavaScript to decode and display them as needed.

This would also remove the downside of the server reliance for HTMX, as we can use the `hx-trigger="load"` to trigger the loading of slides right when they connect, or using the `every` trigger to continue trying to get the slides if the server fails.

The big upside for this approach is just that it makes things slightly easier than a pure java script approach, but it does still share the same downsides as the pure Javascript approach.

For a medium amount of participants this would be an ideal choice. It takes strain away from the developers while also allowing us to optimize as needed.

### Option D: Do Nothing

If we were to just send out slide text as needed using some kind of websocket or SSE, then have it automatically rendered to the screen, it could work with a low number of participants but it's not scalable or maintainable in the long run.

It could very easily cause a server crash or huge lag in receiving lyrics, which would make the product useless as a real-time translated lyric display.

## Analysis

### Performance
| | Positive | Negative |
| :--- | :--- | :--- |
| **HTMX** | On a simulated 3G connection, the gap between HTMX and React initial load times can be 2-3 seconds, with HTMX being significantly faster [(Source)](https://www.softwareseni.com/htmx-vs-react-performance-and-architecture-deep-dive/) | HTMX performance scales at server side, not at client side unlike React which scales with client side [(Source)](https://www.softwareseni.com/htmx-vs-react-performance-and-architecture-deep-dive/)|
| **Pure JS** | JavaScript performance is highly reliant on the developer itself, however when compared to frameworks like React, can be much much faster. [(Source)](https://stackoverflow.com/questions/68115288/react-js-vs-vanilla-js-performance-28x-faster) | Just as before, JavaScript is reliant on the skills of the developer. If the developer makes bad code, it will be inefficient. |
| **Hybrid** | Hybrid model could allow for the exceptional initial load time provided by HTMX (for getting slides content) and the fast performance provided by pure JavaScript | Performance of SSEs and how slide transitions are handled are reliant on the skills of the developer writing the code. |

### Complexity
| | Positive | Negative |
| :--- | :--- | :--- |
| **HTMX** | It hides / removes a lot of the frontend complexity when dealing with client-server interactions. You don't need to write JS for them [(Source)](https://www.erikheemskerk.nl/htmx-simplicity/) | While HTMX presents itself as a simple frontend library, and for the most part it is, HTMX also makes it harder to make forms, handle HTTP errors, test, and manage states [(Source)](https://finance.biggo.com/news/202508120353_HTMX_Reality_Check_Hidden_Complexity) |
| **Pure JS** | JavaScript is inherently more complicated than tools like HTMX, however the developer pool for JS in web development is huge, and almost any problem can be found and solved through resources like stack overflow | JavaScript is just a more complex language in general. JS has on average more nested scopes, more questions on stack overflow, and most code changes per commit on average. [(Source)](https://medium.com/javascript-non-grata/javascript-is-the-most-complex-programming-language-fc68cfb234bb) |
| **Hybrid** | A hybrid model would allow some of the complexity to be removed from JavaScript, making the developers lives easier | JavaScrit would still need to be used for things such as server side events, where its complexity could make it hard to optimize for our environment. |

## Recommendation and Verification Plan

I propose a hybrid approach to the client side frontend. It would allow us to have the efficiency we need as well as removing just a bit of the complexity for the frontend. 

How we will know it worked:
- **Scale:** If at least 250 people can connect to the server at the same time and receive slide transitions
- **Latency:** If the DOM is updated on people's phones within 50-100ms after the slide transition message was sent
- **Accessability:** If the phone's built-in accessibility features can detect the DOM change for accessibility devices

These criteria should be tested throughout the development process, and should be a required CI/CD element once the walking skeleton is complete.

## Constraints

[HTMX is open source](https://github.com/bigskysoftware/htmx) so no license is needed. No other IP or confidentiality boundaries affect this decision as no sensitive data interacts with something other than our server.

## Delegation and Validation

The only AI tool used in the generation of the RFC was the google AI overview for researching the costs and benefits of the approaches laid out above. All sources used from these searches were cited. 

No checks have yet been implemented as the server-side has yet to be developed. 

This problem was thought of before the frontend was developed, so the research of options were all checked by me thoroughly.

## Decision

N/A