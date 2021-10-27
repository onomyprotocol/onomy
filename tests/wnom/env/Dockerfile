FROM onomy/tests-base:latest

WORKDIR /root/home/wnom
# the build is called form the ../ folder so the line copies ../ folder
COPY env env
WORKDIR /root/home/wnom/env
RUN npm install
RUN mkdir logs
RUN touch logs/eth.log
RUN touch logs/orchestrator.log
CMD tail -f logs/eth.log & tail -f logs/orchestrator.log