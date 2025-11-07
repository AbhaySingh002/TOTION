package app


const AsciiArt = ` ______   ______     ______   __     ______     __   __    
/\__  _\ /\  __ \   /\__  _\ /\ \   /\  __ \   /\ "-.\ \   
\/_/\ \/ \ \ \/\ \  \/_/\ \/ \ \ \  \ \ \/\ \  \ \ \-.  \  
   \ \_\  \ \_____\    \ \_\  \ \_\  \ \_____\  \ \_\\"\_\ 
    \/_/   \/_____/     \/_/   \/_/   \/_____/   \/_/ \/_/ 
                                                           `

const GeneralHelp = "Ctrl+N: New Note • Ctrl+L: List all Notes • Esc: Return to home • Ctrl+C: Quit Totion "
const SaveHelp = "Ctrl+N: New Note • Ctrl+L: List all Notes • Esc: Return to home • Ctrl+S: Save Note • Ctrl+C: Quit Totion"
const ListHelp = "Ctrl+N: New Note • Esc: Return to home • Ctrl+C: Quit Totion • Delete / Backspace: Delete Note • Enter: Open Note"
const SystemPrompt = `"You are an intelligent note assistant that helps users thoughtfully continue their notes.
Continue the note in a natural, meaningful, and concise way — capturing the same tone or emotion.
Do not repeat the existing text. Do not add any labels like "Completion:" or quotes. Make sure that sentence is complete. Don't end or start with the "..." .
Don't add the extra things except the Continuation part, becuse the response is going to be directly implemented in the notes.
Note:
%s
Continuation:"`
const GenaiModel = "gemini-2.5-flash-lite"
const Api_key = "GEMINI-API-KEY"