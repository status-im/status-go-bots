FROM ubuntu:16.04

# setup libraries
RUN apt-get -y update
RUN apt-get -y install openjdk-8-jdk
RUN apt-get install -y wget unzip build-essential file

# setup pathes env
ENV JAVA_HOME /usr/lib/jvm/java-8-openjdk-amd64
ENV ANDROID_HOME /opt/android-sdk-linux
ENV NDK_ROOT /opt/android-ndk-r16b
ENV ANDROID_NDK_HOME ${NDK_ROOT}

# setup Android SDK
RUN cd /opt \
    && wget -q https://dl.google.com/android/repository/sdk-tools-linux-3859397.zip -O android-sdk-tools.zip \
    && unzip -q android-sdk-tools.zip -d ${ANDROID_HOME} \
    && rm android-sdk-tools.zip
ENV PATH ${PATH}:${ANDROID_HOME}/tools:${ANDROID_HOME}/tools/bin:${ANDROID_HOME}/platform-tools:${NDK_ROOT}
RUN yes | sdkmanager --licenses
RUN sdkmanager "emulator" "tools" "platform-tools"

# setup Android SDKmanager
RUN yes | sdkmanager \
    "platforms;android-27" \
    "platforms;android-26" \
    "platforms;android-25" \
    "platforms;android-24" \
    "platforms;android-23" \
    "platforms;android-22" \
    "platforms;android-21" \
    "platforms;android-19" \
    "platforms;android-17" \
    "platforms;android-15" \
    "build-tools;27.0.3" \
    "build-tools;27.0.2" \
    "build-tools;27.0.1" \
    "build-tools;27.0.0" \
    "build-tools;26.0.2" \
    "build-tools;26.0.1" \
    "build-tools;25.0.3" \
    "build-tools;24.0.3" \
    "build-tools;23.0.3" \
    "extras;android;m2repository" \
    "extras;google;m2repository" \
    "extras;google;google_play_services"

# setup maven
RUN apt-get purge maven maven2 \
    && apt-get update \
    && apt-get -y install maven \
    && mvn --version

# setup ndk
RUN wget -q https://dl.google.com/android/repository/android-ndk-r16b-linux-x86_64.zip \
    && unzip -q android-ndk-r16b-linux-x86_64.zip -d /opt \
    && rm android-ndk-r16b-linux-x86_64.zip

RUN apt-get install -y software-properties-common
RUN add-apt-repository ppa:longsleep/golang-backports
RUN apt-get -y update
RUN apt-get -y install golang-go
ENV PATH="${PATH}:/root/go/bin"

RUN apt-get install -y git
RUN apt-get install -y curl



WORKDIR /root/go/src/github.com/status-im/status-go

CMD ["/bin/bash"]



