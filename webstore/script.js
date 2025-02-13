let conversations; // Stores the conversations
let currentChat; // Points to the active conversation

const startSpeach = document.getElementById("speechButton");
const messageInput = document.getElementById("messageInput");
const languageSelect = document.getElementById("languageSelector");

const SpeechRecognition = window.SpeechRecognition || window.webkitSpeechRecognition; // Used in converting speech to text
const synth = window.speechSynthesis; // Used in converting text to speech
let voices = synth.getVoices(); // Stores the voices used in text to speech synth

document.getElementById("input_server_ip").value = window.location.host;

if (SpeechRecognition) {

    const recognition = new SpeechRecognition();

    recognition.continuous = false; // Stop when speech ends
    recognition.lang = 'en-US'; // Language
    recognition.interimResults = true; // Show interim results

    // Start recognition
    startSpeach.addEventListener('click', () => {
    recognition.start();
    });

    languageSelect.addEventListener('change', () => {
        recognition.lang = languageSelect.value;
    });

    // Capture the result
    recognition.onresult = (event) => {
        let transcript = '';
        for (let i = 0; i < event.results.length; i++) {
            transcript += event.results[i][0].transcript;
        }
        messageInput.value = transcript;
        //text = transcript;
    };

    // Actions to do when the capture ends
    recognition.onend = (event) => {
        MessageLLM();
    }

    // Error handling
    recognition.onerror = (event) => {
        messageInput.value = 'Error occurred: ' + event.error;
    };
} else {
    startSpeach.style.backgroundColor = "grey";
};


function populateVoiceList() {
    const voiceSelect = document.getElementById("voiceSelector");

    voices.forEach(voice => {
        const option = document.createElement("option");
        option.textContent = `${voice.name} (${voice.lang})`;

        if (voice.default) {
            option.textContent += " — DEFAULT";
        }

        option.setAttribute("data-lang", voice.lang);
        option.setAttribute("data-name", voice.name);
        voiceSelect.appendChild(option);
    });
};

window.speechSynthesis.onvoiceschanged = function() {
    voices = synth.getVoices();
    populateVoiceList();
};  

// Function to switch chats
function switchChat(index) {
    currentChat = index;

    // Update active class
    document.querySelectorAll(".chat-list li").forEach((li, i) => {
        if (li.innerHTML == `Chat-${currentChat}`) {
            li.classList.add("active");
        } else {
            li.classList.remove("active");
        }
    });

    document.querySelector(".chat-header").textContent = currentChat;
    loadMessages();
};

// Function to load messages
function loadMessages() {
    const chatBox = document.querySelector(".chat-box");
    chatBox.innerHTML = "";

    conversations.get(currentChat).forEach(msg => {
        const messageDiv = document.createElement("div");
        messageDiv.classList.add("message", msg.type);
        messageDiv.textContent = msg.text;
        chatBox.appendChild(messageDiv);
    });

    chatBox.scrollTop = chatBox.scrollHeight;
};

function createConversation() {
    setCurrentChat()
    var chatlist = document.getElementsByClassName("chat-list")[0];
    var conversation = document.createElement("li");
    conversation.setAttribute('onclick', `switchChat('${currentChat}')`);
    conversation.innerHTML = `Chat-${currentChat}`;
    chatlist.appendChild(conversation);
    conversations.set(currentChat, [{ text: "New conversation. Send a message to begin the conversation", type: "received" }]);
    document.querySelector(".chat-header").textContent = currentChat;
};

function MessageLLM() {

    url = "http://" + `${document.getElementById("input_server_ip").value}`+ "/api/chat"

    const input = document.getElementById("messageInput");
    const text = input.value.trim();
    if (text === "") return;

    // Add message to conversation
    conversations.get(currentChat).push({ text, type: "sent" });
    loadMessages();

    // Clear input
    input.value = "";

    body = JSON.stringify({
            userid: currentChat,
            prompt: text
        });

    const response = fetch(url, {
        method: 'POST',
        body: body,
         headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
        },
    })
    .then(response => response.text())
    .then(data => { 
        setTimeout(() => {
            conversations.get(currentChat).push({ text: data, type: "received" });
            loadMessages();

            // If a voice is selected then speak out the LLM response
            if (document.getElementById("voiceSelector").value !== "") {
                const response = new SpeechSynthesisUtterance(data);
                const selectedOption = document.getElementById("voiceSelector").selectedOptions[0].getAttribute("data-name");

                voices.forEach(voice => { 
                    if (voice.name === selectedOption) {
                        response.voice = voice;
                    }
                });

                response.pitch = 1;
                response.rate = 1;
                synth.speak(response);
            }
         }, 1000);
    });
 };

function setCurrentChat() {
    currentChat = (`${Math.floor(Math.random() * 10000)}`).toString();
};

function setupMap() {
    conversations = new Map();
};



setupMap();
setCurrentChat();
createConversation(currentChat);
loadMessages();

document.querySelector("#messageInput").addEventListener("keyup", event => {
    if(event.key !== "Enter") return; // Use `.key` instead.
    MessageLLM() // Things you want to do.
    event.preventDefault(); // No need to `return false;`.
});


