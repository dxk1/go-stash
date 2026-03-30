FROM csighub.tencentyun.com/scf/component:scf_baseimage_v11

# 支持多架构构建
ARG TARGETPLATFORM
ARG BUILDPLATFORM
RUN echo "I am running on $BUILDPLATFORM, building for $TARGETPLATFORM"

WORKDIR /usr/local/services/scf_stash/

COPY stash /usr/local/services/scf_stash/
# 复制配置文件到容器中
COPY etc/config.yaml /usr/local/services/scf_stash/scf_stash.yaml
# 复制健康检查脚本
COPY healthcheck.sh /usr/local/services/scf_stash/

# 确保二进制文件和脚本具有执行权限
RUN chmod +x /usr/local/services/scf_stash/stash && \
    chmod +x /usr/local/services/scf_stash/healthcheck.sh

# 添加健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
  CMD /usr/local/services/scf_stash/healthcheck.sh

CMD ["/usr/local/services/scf_stash/stash"]