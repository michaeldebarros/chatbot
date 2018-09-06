# chatbot

This is a simple chatbot, initially intended to help the user book their vacation on another system.

The frontend is made with [chat-bubble](https://github.com/dmitrizzle/chat-bubble).  I had to tweek it a little bit because in this project the conversation is all held in a JS object. But my project was designed to work with a Go backend that serves and receives JSON to and from chat-bubble.

The conversation is based on a struct named passosFerias (vacation steps in portuguese).
