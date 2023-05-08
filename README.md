# Gptx

[![test](https://github.com/kohkimakimoto/gptx/actions/workflows/test.yml/badge.svg)](https://github.com/kohkimakimoto/gptx/actions/workflows/test.yml)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/kohkimakimoto/gptx/blob/main/LICENSE)
[![release](https://img.shields.io/github/release/kohkimakimoto/gptx/all.svg)](https://github.com/kohkimakimoto/gptx/releases)

An extensible command-line utility powered by ChatGPT, designed to enhance your productivity.

## Overview

Gptx is designed to bridge the gap between the [ChatGPT API](https://platform.openai.com/docs/api-reference/chat) and your shell environment.
Seamlessly integrate ChatGPT's natural language processing capabilities with familiar shell functionalities.
With Gptx, you can easily script, automate, and elevate your shell tasks with the power of ChatGPT, unlocking a new world of AI-assisted productivity in a terminal environment.

### Features

- [Chat](#simple-chat-messages) with ChatGPT from your terminal
- Manage a series of messages as a [conversation](#conversations)
- [Cache](#cache) responses from ChatGPT API
- Highly customizable with [Hooks](#hooks) and [Custom subcommands](#custom-subcommands)

## Installation

### Precompiled binaries

Gptx is a single binary command line program. You can download a precompiled binary at the Github releases page.

[Download the latest version](https://github.com/kohkimakimoto/gptx/releases/latest)

### Homebrew

```sh
brew install kohkimakimoto/gptx/gptx 
```

### From source

```sh
go install github.com/kohkimakimoto/gptx/cmd/gptx@latest
```

## Usage

After installing, you can run `gptx init` command to initialize the configuration.

```sh
gptx init
```

This command creates a directory `~/.gptx`. It has a configuration file `config.toml`.

> :information_source: Note: You can also change the directory path by setting the environment variable `GPTX_HOME`.
> This directory is used to store the configuration file, along with other files utilized by the gptx command.

Open the configuration file `~/.gptx/config.toml` with your favorite editor, and enter your OpenAI API key in the `openai_api_key` field.

```toml
# ~/.gptx/config.toml
openai_api_key = "sk-*******"
```

### Simple chat messages

You can chat with ChatGPT by running the `gptx chat` or `gptx c` command.

```sh
gptx chat "What is the capital city of Japan?"
# -> The capital city of Japan is Tokyo.
````

Gptx can accept prompt text from STDIN.

```sh
echo "What is the most famous landmark in Tokyo?" | gptx chat
# -> The most famous landmark in Tokyo is the Tokyo Tower.
```

> :information_source: Note: You can use both STDIN and arguments at the same time. They are concatenated with a newline such as `STDIN + "\n" + ARGUMENTS`.

https://user-images.githubusercontent.com/761462/235860236-7704caa2-6f5f-49a2-b7f8-b472ec255e15.mp4

### Conversations

Conversations in Gptx consist of a series of messages. Each conversation has one or more messages.
Messages comprise user's message or assistant's message. A user's message is a text that you send to ChatGPT, and a assistant's message is a text that ChatGPT returns to you.

By default, Gptx creates a new conversation for each chat process described above. You can list the created conversations by running the `gptx list` or `gptx ls` command.

```sh
gptx list
```

````
ID   PROMPT                                         MESSAGES   NAME   LABEL   HOOKS   CREATED                ELAPSED
 1   What is the capital city of Japan?                    2                          2023-05-03T06:40:24Z   21 seconds ago
 2   What is the most famous landmark in Tokyo?            2                          2023-05-03T06:40:37Z   8 seconds ago
````

> :information_source: Note: Conversations are saved to the internal database file `$GPTX_HOME/gptx.db`.

You can display the conversation details by running the `gptx inspect` or `gptx i` command with conversation ID.

```sh
gptx inspect -p 1
```

> :information_source: Note: The -p option is used for pretty-printing JSON data.

```json
{
  "id": 1,
  "prompt": "What is the capital city of Japan?",
  "created_at": "2023-05-03T06:40:24.465656Z",
  "messages": [
    {
      "role": "user",
      "content": "What is the capital city of Japan?"
    },
    {
      "role": "assistant",
      "content": "The capital city of Japan is Tokyo."
    }
  ]
}
```

If you want to send a message in an existing conversation context, you can use the `--resume` or `-r` option.

```sh
gptx chat --resume 1 "What about the USA?"
# -> The capital city of the United States of America (USA) is Washington D.C.
```

You can assign a unique name to a conversation using the `--name` or `-n` option. This is helpful if you want to resume the conversation later.

```sh
# Create a new conversation with a name.
gptx chat -n city "What is the capital city of Japan?"
# -> The capital city of Japan is Tokyo.

# Resume the conversation by name.
gptx chat -r city "What about the USA?"
# -> The capital city of the United States of America (USA) is Washington D.C.
```

https://user-images.githubusercontent.com/761462/235863838-e1792bdb-542f-426e-8dba-bd62b1d655c4.mp4

### Cache

By default, Gptx caches the response from ChatGPT API. When you send the exact same message to ChatGPT API, Gptx returns the cached response instead of sending a request to ChatGPT API.
If you don't want to use the cache, you can use the `--no-cache` option.

```sh
gptx chat --no-cache "What is the capital city of Japan?"
```

You can also clear the all cache by running the `gptx clean` command.

```sh
gptx clean
```

### Interactive mode

Gptx has an interactive (REPL) mode. You can enter the interactive mode by running the `gptx chat` command with `--interactive` or `-i` option.

```sh
gptx chat -i
```

```
Welcome to gptx v0.0.1.
Type ".help" for more information.
> What is the capital city of Japan?
The capital city of Japan is Tokyo.
> What about the USA?
The capital city of the United States of America (USA) is Washington D.C.
>
```

> :information_source: Note: In the interactive mode, your input messages are in the same conversation context.

https://user-images.githubusercontent.com/761462/235866177-eb76ca9c-3f81-406e-966c-a196899ae282.mp4

## Configuration

The configuration file must be written in [TOML](https://github.com/toml-lang/toml).
By default, Gptx loads the configuration file from  `~/.gptx/config.toml`.
If you want to change the path, you can set the environment variable `GPTX_HOME` and then it will load the configuration file from `$GPTX_HOME/config.toml`.

### Example

```toml
# OpenAI API Key. You can override this value by using the OPENAI_API_KEY environment variable.
openai_api_key = ""

# Default model for Chat API
model = "gpt-3.5-turbo"

# Maximum number of cached responses.
max_cache_length = 100
```

## Hooks

Gptx hooks offer a powerful mechanism for extending the functionality of your Gptx processes.
Designed to inspect and filter the messages, hooks provide the flexibility to execute custom logic both before and after submitting a request to the ChatGPT API.
Hooks can generate dynamic prompts, parse outputs, and implement any other custom behavior to tailor your ChatGPT experience to your specific needs.

### Example: Shell Hook

Gptx has a reference implementation of the hook called *Shell Hook*. You can use it by running the `gptx chat` command with `--hook` or `-H` option. Try the following command.

```sh
gptx chat --hook shell "start http server with python3"
```

You will see the following output.

```
python3 -m http.server
Choose an action: Run (r), Copy to clipboard (c), or Nothing (n):
```

Using the *Shell Hook*, Gptx is able to generate shell commands based on your input message and then asks for additional action.

https://user-images.githubusercontent.com/761462/236078907-6016b8e5-83cd-44ed-bd68-5467838678bd.mp4

### How do hooks actually work?

Hooks are command line programs that you can write in any programming language.
The Gptx chat process goes through a series of steps, with hooks being executed along the way to provide opportunities for modifying the chat process behavior.

The following diagram illustrates the chat process with hooks incorporated.

![hooks-diagram](https://user-images.githubusercontent.com/761462/236674975-1d0eab37-8a6a-4904-a853-acf8d2e90636.svg)

As you can see in the diagram, there are various points at which hooks are executed: `pre-message`, `post-message`, and `finish`.
These execution points are called [Types of hooks](#types-of-hooks).
Hooks specified by the `--hook` or `-H` option are executed at all of these points.
You can retrieve the hook execution point from the `GPTX_HOOK_TYPE` environment variable in the hook programs.

So, the minimum hook written in Bash script is like the following:

```bash
#!/usr/bin/env bash
set -e -o pipefail

case "$GPTX_HOOK_TYPE" in
  'pre-message')
    echo "run pre-message hook"
    ;;
  'post-message')
    echo "run post-message hook"
    ;;
  'finish')
    echo "run finish hook"
    ;;
  *)
    echo "invalid hook kind: $GPTX_HOOK_TYPE" 1>&2
    exit 1
    ;;
esac
```

Hooks must have a `gptx-hook-` prefix in their file name and be executable.
For example, if you want to use the above minimum example hook, you should create a file named `gptx-hook-example` and make it executable.

You can place the hooks in the directory indicated by the PATH environment variable, just like general command line programs.
Gptx also provides a `libexec` directory under the Gptx home directory (`~/.gptx/libexec` by default) for placing executable programs related to Gptx.
The `libexec` directory is automatically added to the PATH environment variable while Gptx is running.
It is recommended to place the hook programs in the `libexec` directory.

> :information_source: Note: [*Shell Hook*](#example-shell-hook) is also placed in the `libexec` directory.

If you run the above example hook, you will see the following output.

```sh
gptx chat -H example "Say hello"
# -> run pre-message hook
# -> run post-message hook
# -> Hello! How may I assist you today?
# -> run finish hook
```

That's it! These are the requirements for implementing the minimum hook.
If you want more details, please read the following sections.

### Environment variables

Gptx uses environment variables to pass information to hooks. The following environment variables are available to all hooks.

- `GPTX_HOOK_TYPE`: The type of hook in which the hook is executed. The value is one of `pre-message`, `post-message`, or `finish`.
- `GPTX_MESSAGE_INDEX`: The index of the current message within the conversation. The value is an integer starting from `0`, with `0` representing the first message.
- `GPTX_CONVERSATION_ID`: The ID of the conversation. If the hook processes a new conversation and is in the `pre-message` stage, the value is `0`. This indicates that the conversation has not been saved and does not have an ID yet.

### Types of hooks

Hooks are executed at various points during the chat process.
You need to implement appropriate logic in the hook programs based on the execution points.
If no action is needed at a specific point, you should write programs that ignore it.

The following sections describe the types of hooks in detail.

#### pre-message

This hook is executed before sending a message to the ChatGPT API, allowing you to modify the message in advance.
You can read the input prompt text from the file whose path is specified by the `GPTX_PROMPT_FILE` environment variable.
The `GPTX_PROMPT_FILE` environment variable also indicates the path to the file where the hook should write the modified prompt text.

Example: `gptx-hook-translation`

```bash
#!/usr/bin/env bash
set -e -o pipefail

case "$GPTX_HOOK_TYPE" in
  'pre-message')
    # read the input prompt text
    prompt=$(cat "$GPTX_PROMPT_FILE")
    # modify the prompt text and write it to the file
    cat << EOF > "$GPTX_PROMPT_FILE"
Translate the following text into Japanese.
Text: $prompt
Japanese:
EOF
    ;;
  'post-message')
    # for post-message code...
    ;;
  'finish')
    # for post-message code...
    ;;
  *)
    echo "invalid hook kind: $GPTX_HOOK_TYPE" 1>&2
    exit 1
    ;;
esac
```

Output:

```sh
gptx chat -H translation hello
# -> こんにちは
```

#### post-message

This hook is executed after receiving a message from the ChatGPT API, allowing you to modify the AI generated output.
You can read the AI-generated output from the file whose path is specified by the `GPTX_COMPLETION_FILE` environment variable.
The `GPTX_COMPLETION_FILE` environment variable also indicates the path to the file where the hook should write the modified AI-generated output.

#### finish

This hook is executed at the last stage of the chat process.
You can read the final output of the chat process from the file whose path is specified by the `GPTX_COMPLETION_FILE` environment variable.
You can use this hook to run additional processing with the chat process output.

### Caveats

- You can assign multiple hooks to the `--hook` or `-H` option. In this case, the hooks are executed in the order in which they are specified.
- You can only assign hooks to the new conversation. Assigned hooks are saved in the conversation object. If you [resume](#conversations) the conversation, the saved hooks will be executed.
- The hook can use an exit code `3` as a *Cancel* signal. If the hook returns an exit code `3`, the Gptx process is terminated immediately with no error.

## Custom subcommands

Gptx allows you to add custom subcommands.
Gptx recognizes commands with the `gptx-` prefix as these custom subcommands.
For example, if you create a command named `gptx-foo`, you can run it as `gptx foo`.

As I mentioned in the [How do hooks actually work?](#how-do-hooks-actually-work) section, you can place the executable programs in the `libexec` directory under the Gptx home directory (`~/.gptx/libexec` by default).
The custom subcommands should also be placed in the libexec directory, as it is the recommended location.

## Author

Kohki Makimoto <kohki.makimoto@gmail.com>

## License

The MIT License (MIT)
