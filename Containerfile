FROM registry.access.redhat.com/ubi8/ubi
RUN dnf -y install java-11-openjdk-devel && \
    echo 'public class ZonedHello { public static void main(String[] args) { System.out.println(java.time.ZonedDateTime.now()); } }' > ZonedHello.java && \
    javac ZonedHello.java

CMD ["/bin/bash", "-c", "ls -la /etc/localtime && mount | grep zoneinfo || true && echo \"Now with default timezone:\" && date && echo \"Java default sees the following timezone:\" && java ZonedHello && echo \"Forcing UTC:\" && TZ=Etc/UTC date"]
