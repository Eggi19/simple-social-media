# simple-social-media

## Features

- **User Authentication**: 
  - Secure user login and registration.
  - Passwords are securely stored, ensuring user data safety.

- **Posting, Liking, and Commenting**:
  - Users can create and post tweets.
  - Liking functionality allows users to like tweets.
  - Users can comment on tweets.

- **Notifications**:
  - Real-time notifications inform users when someone likes or comments on their tweet and when other user follows the user.
  - Notifications are managed through Firebase Cloud Messaging.

- **Follow/Unfollow**:
  - Users can follow or unfollow each other.
  - Follower and following count and information are stored using firestore.

  ## Tech Stack

- **Golang**: Backend programming language for handling business logic and API requests.
- **Firebase**:
  - **Firestore**: For structured data storage of user profiles, followers and following count and data.
  - **Realtime Database**: For real-time tweet's like.
  - **Firebase Cloud Messaging (FCM)**: To send push notifications when a comment is made on a post, user follows other user and when liking the tweet.

  ## How to Run the Project

1. **Set Up Firebase**:
   - Create a Firebase project.
   - Enable Firestore, and Realtime Database in the Firebase Console.
   - Download the firebase service account key file and place it in the project directory as `firebase-service-account-key.json`, you can see the example in directory.

2. **Set Up .env**
   - You can see the .env structure in .env.example file

3. **Create Your Postgres Database**
   - Create your new database
   - run ddl.sql on your new database (you can find the file in sql folder)

4. **Install Dependencies**: