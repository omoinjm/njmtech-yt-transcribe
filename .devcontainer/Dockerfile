FROM mcr.microsoft.com/devcontainers/base:ubuntu

ARG USERNAME=vscode

# [Optional] Uncomment this section to install additional OS packages.
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends software-properties-common \
    ffmpeg \
    git \
    build-essential \
    cmake \
    pkg-config \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Set up bash history
RUN SNIPPET="export PROMPT_COMMAND='history -a' && export HISTFILE=/commandhistory/.bash_history" \
    && mkdir /commandhistory \
    && touch /commandhistory/.bash_history \
    && chown -R $USERNAME /commandhistory \
    && echo "$SNIPPET" >> "/home/$USERNAME/.bashrc"

# Install yt-dlp
RUN sudo curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp \
    && sudo chmod a+rx /usr/local/bin/yt-dlp

# Install whisper.cpp
RUN git clone https://github.com/ggml-org/whisper.cpp.git /whisper.cpp \
    && cd /whisper.cpp \
    && sh ./models/download-ggml-model.sh base.en \
    && cmake -B build -DCMAKE_BUILD_TYPE=Release \
    && cmake --build build -j \
    && cp build/bin/whisper-cli /usr/local/bin/whisper-cli
