FROM csighub.tencentyun.com/scf/component:scf_baseimage_v11

# 支持多架构构建
ARG TARGETPLATFORM
ARG BUILDPLATFORM
RUN echo "I am running on $BUILDPLATFORM, building for $TARGETPLATFORM"

WORKDIR /usr/local/services/scf_stash/

COPY stash /usr/local/services/scf_stash/

# 确保二进制文件具有执行权限
RUN chmod +x /usr/local/services/scf_stash/stash

CMD ["/usr/local/services/scf_stash/stash"]