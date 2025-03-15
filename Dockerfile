# Use Go 1.24 as the builder stage
FROM golang:1.24 AS builder

# Install required X11 development libraries for Ebiten
RUN apt-get update && apt-get install -y \
    libgl1-mesa-dev xorg-dev libx11-dev

# Set working directory
WORKDIR /app

# Copy and build the game
COPY . .
RUN go mod tidy && go build -o snake

# Create the final image
FROM ubuntu:22.04

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    xvfb x11vnc fluxbox libgl1-mesa-glx libgl1-mesa-dri \
    libxcursor1 libxi6 libxrandr2 && \
    rm -rf /var/lib/apt/lists/*

# Set environment variables for headless display
ENV DISPLAY=:1
ENV SCREEN_WIDTH=400
ENV SCREEN_HEIGHT=400
ENV SCREEN_DEPTH=24

# Copy the built game
COPY --from=builder /app/snake /usr/local/bin/snake

# Create script to launch X server and game
RUN echo '#!/bin/bash\n\
Xvfb :1 -screen 0 ${SCREEN_WIDTH}x${SCREEN_HEIGHT}x${SCREEN_DEPTH} &\n\
sleep 2\n\
fluxbox &\n\
x11vnc -display :1 -forever -nopw -rfbport 15900 &\n\
sleep 2\n\
/usr/local/bin/snake\n' > /start.sh && chmod +x /start.sh

# Start VNC when the container runs
CMD ["/start.sh"]
