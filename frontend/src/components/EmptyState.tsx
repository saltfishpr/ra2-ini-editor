import { theme, Typography } from "antd";

interface EmptyStateProps {}

const EmptyState = ({}: EmptyStateProps) => {
  const { token } = theme.useToken();

  return (
    <div
      style={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        height: "100%",
        flexDirection: "column",
      }}
    >
      <Typography.Title level={4} style={{ color: token.colorTextSecondary }}>
        请选择一个单位查看详情
      </Typography.Title>
      <Typography.Text type="secondary">
        从左侧列表中选择一个单位来查看其详细信息
      </Typography.Text>
    </div>
  );
};

export default EmptyState;
