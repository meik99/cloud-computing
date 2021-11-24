# Proposal - Video Streaming Platform

The idea is to create a web platform for video streaming and sharing.
Like YouTube, it should be possible to upload and then view videos.
These videos can then be shared with other people or be set to public.

## Potential Data and Services

* Metadata for Videos
  * Id
  * Title
  * Description
  * Tags
  * Storage path
* Access of Videos
  * Owner
  * Shared With
  * Is Public
* Videos
  * Handles streaming of a video
* User
  * Id
  * Name
* Profile for a User
  * Owner
  * Description
* Rating of Videos
  * Likes
  * Dislikes
  * Video Id
* Comments
  * User Id
  * Text
  * Replies

## Data Storage

A MongoDB instance could be used to store the metadata, since it is already made to be distributable. 
The videos themselves should be stored as file, accessed by a custom service. 
Alternatively to a custom service, research could be done to find already existing solutions.

## [Operator](https://sdk.operatorframework.io/)

If time allows, an operator could also be written and deployed. 
This operator can make sure, that the services, i.e. as pods, run correctly or automatically deploy them as needed.
It would make sure the infrastructure is set up in a correct way, or report any errors if not.

## Responsibilities

From the above services, the following responsibilities are assigned.

[Daniel Klösler](https://github.com/Ethlaron):
* Video Metadata Service
* Video Access Service (Authorization)

[Christoph Pargfrieder](https://github.com/ChristophPargfrieder):
* Comments Service
* Rating Service
* Video Service

[Michael Rynkiewicz](https://github.com/meik99): 
* Frontend
* User Service (Authentication)
* Profile Service
* Operator

## Milestones

The following milestones are defined as use-case scenarios.

### Milestone 1

It is possible to authenticate a Google account against the Google OAuth 2 endpoint and show their name on the frontend.
Authentication is optional, i.e. the frontend does not try to authenticate automatically.
Only the landing page can be viewed unauthenticated.
An authenticated user can view, edit and delete their own profile.

### Milestone 2

Videos can be uploaded and viewed.
The frontend provides a workflow to correctly upload a video file and user defined metadata.
A viewing showing a profile lists uploaded videos.
From there it is possible to watch a video and view its metadata.

### Milestone 3

Videos can be viewed and commented by users other than the original poster.
Access rights for a video can be defined by the poster.
Those access rights may be, for example, public, restricted and private.
Additionally, videos can be rated and searched for.
