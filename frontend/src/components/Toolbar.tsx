import { DeleteOutlined, PlusOutlined, SaveOutlined } from "@ant-design/icons";
import { Button, Form, Input, Modal, Select, Space } from "antd";
import { useState } from "react";

interface ToolbarProps {
  onAdd: (type: string, name: string) => void;
  onSave: () => void;
  onDelete: () => void;
  canSave: boolean;
  canDelete: boolean;
}

const Toolbar: React.FC<ToolbarProps> = ({
  onAdd,
  onSave,
  onDelete,
  canSave,
  canDelete,
}) => {
  const [modalOpened, setModalOpened] = useState(false);
  const [unitType, setUnitType] = useState("");
  const [unitName, setUnitName] = useState("");

  return (
    <Space className="toolbar" size="small">
      <Button
        type="primary"
        icon={<PlusOutlined />}
        onClick={() => setModalOpened(true)}
      >
        新增
      </Button>
      <Modal
        title="新增单位"
        open={modalOpened}
        onCancel={() => setModalOpened(false)}
        onOk={() => {
          onAdd(unitType, unitName);
          setModalOpened(false);
          setUnitType("");
          setUnitName("");
        }}
        width={400}
      >
        <Form layout="vertical">
          <Form.Item label="单位类型" required>
            <Select
              style={{ width: "100%" }}
              placeholder="请选择单位类型"
              options={[
                { value: "infantry", label: "步兵" },
                { value: "vehicle", label: "战车" },
                { value: "aircraft", label: "飞行物" },
                { value: "building", label: "建筑" },
              ]}
              value={unitType}
              onSelect={(value) => setUnitType(value)}
            />
          </Form.Item>
          <Form.Item label="单位名称" required>
            <Input
              placeholder="请输入单位名称，如：E1"
              value={unitName}
              onChange={(e) => setUnitName(e.target.value)}
            />
          </Form.Item>
        </Form>
      </Modal>
      <Button icon={<SaveOutlined />} onClick={onSave} disabled={!canSave}>
        保存
      </Button>
      <Button
        danger
        icon={<DeleteOutlined />}
        onClick={onDelete}
        disabled={!canDelete}
      >
        删除
      </Button>
    </Space>
  );
};

export default Toolbar;
