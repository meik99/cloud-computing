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

## Data Storage

A MongoDB instance could be used to store the metadata, since it is already made to
be distributable. The videos themselves should be stored as file, 
accessed by a custom service. Alternatively to a custom service, research could
be done to find already existing solutions.

## Operator

If time allows, an operator could also be written and deployed. This operator can make
sure, that the services, i.e. as pods, run correctly or automatically deploy them as needed.
It would make sure the infrastructure is set up in a correct way, or report any errors if not.
