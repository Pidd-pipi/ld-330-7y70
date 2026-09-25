import { HistoryOutlined, PlusOutlined, SearchOutlined, SwapOutlined } from '@ant-design/icons';
import {
  Alert,
  Button,
  Card,
  Descriptions,
  Drawer,
  Form,
  Input,
  InputNumber,
  Modal,
  Popconfirm,
  Radio,
  Select,
  Space,
  Table,
  Tag,
  Typography,
  message,
} from 'antd';
import { useEffect, useState } from 'react';
import { api } from '../api/client';
import { useAuth } from '../store/auth';
import type { Patient, PatientMergeLog, PatientMergePreview } from '../types';

const genderOptions = ['男', '女', '其他'];

export default function Patients() {
  const { user } = useAuth();
  const isAdmin = user?.role === 'admin';
  const [items, setItems] = useState<Patient[]>([]);
  const [open, setOpen] = useState(false);
  const [mergeOpen, setMergeOpen] = useState(false);
  const [logsOpen, setLogsOpen] = useState(false);
  const [logs, setLogs] = useState<PatientMergeLog[]>([]);
  const [q, setQ] = useState('');
  const [keepPatientId, setKeepPatientId] = useState<number>();
  const [mergedPatientId, setMergedPatientId] = useState<number>();
  const [reason, setReason] = useState('');
  const [preview, setPreview] = useState<PatientMergePreview | null>(null);
  const [previewLoading, setPreviewLoading] = useState(false);
  const [mergeLoading, setMergeLoading] = useState(false);

  const load = async () => {
    try {
      const r = await api.get<{ items: Patient[] }>('/patients', { keyword: q, page_size: 100 });
      setItems(r.data.items);
    } catch (e) {
      message.error((e as Error).message);
    }
  };

  useEffect(() => {
    void load();
  }, []);

  const save = async (v: object) => {
    try {
      await api.post('/patients', v);
      message.success('患者档案已建立');
      setOpen(false);
      void load();
    } catch (e) {
      message.error((e as Error).message);
    }
  };

  const openMerge = (keepId?: number) => {
    setKeepPatientId(keepId);
    setMergedPatientId(undefined);
    setReason('');
    setPreview(null);
    setMergeOpen(true);
  };

  const loadPreview = async () => {
    if (!keepPatientId || !mergedPatientId || keepPatientId === mergedPatientId || reason.trim().length < 2) {
      message.warning('请选择两套不同档案，并填写至少 2 个字的合并原因。');
      return;
    }
    setPreviewLoading(true);
    try {
      const r = await api.post<PatientMergePreview>('/patients/merge-preview', {
        keep_patient_id: keepPatientId,
        merged_patient_id: mergedPatientId,
        reason,
      });
      setPreview(r.data);
    } catch (e) {
      setPreview(null);
      message.error((e as Error).message);
    } finally {
      setPreviewLoading(false);
    }
  };

  const confirmMerge = async () => {
    if (!preview) return;
    setMergeLoading(true);
    try {
      await api.post('/patients/merge', {
        keep_patient_id: keepPatientId,
        merged_patient_id: mergedPatientId,
        reason,
      });
      message.success('档案已合并，病历和处方均已转入保留档案。');
      setMergeOpen(false);
      void load();
    } catch (e) {
      message.error((e as Error).message);
    } finally {
      setMergeLoading(false);
    }
  };

  const loadLogs = async () => {
    setLogsOpen(true);
    try {
      const r = await api.get<PatientMergeLog[]>('/patient-merge-logs');
      setLogs(r.data);
    } catch (e) {
      message.error((e as Error).message);
    }
  };

  const patientOptions = items.map((x) => ({
    value: x.id,
    label: `${x.name} · ${x.record_no} · ${x.id_card}`,
  }));

  return (
    <Card
      className="page-card"
      title={<Typography.Title level={3} className="page-title">患者档案</Typography.Title>}
      extra={
        <Space>
          <Input
            placeholder="姓名 / 身份证 / 手机号"
            value={q}
            onChange={(e) => setQ(e.target.value)}
            onPressEnter={() => void load()}
            suffix={<SearchOutlined />}
          />
          <Button onClick={() => void load()}>检索</Button>
          {isAdmin && <Button icon={<HistoryOutlined />} onClick={() => void loadLogs()}>合并记录</Button>}
          {isAdmin && <Button icon={<SwapOutlined />} onClick={() => openMerge()}>发起档案合并</Button>}
          <Button type="primary" icon={<PlusOutlined />} onClick={() => setOpen(true)}>新建档案</Button>
        </Space>
      }
    >
      <Table
        rowKey="id"
        dataSource={items}
        columns={[
          { title: '档案编号', dataIndex: 'record_no' },
          { title: '姓名', dataIndex: 'name' },
          { title: '性别', dataIndex: 'gender', render: (v) => <Tag color={v === '男' ? 'blue' : 'magenta'}>{v}</Tag> },
          { title: '年龄', dataIndex: 'age' },
          { title: '当前病历数', render: (_, x) => <Tag color="blue">{x.record_count || 0}</Tag> },
          { title: '身份证号', dataIndex: 'id_card' },
          { title: '联系方式', dataIndex: 'phone' },
          { title: '过敏史', dataIndex: 'allergies', ellipsis: true },
          { title: '既往病史', dataIndex: 'medical_history', ellipsis: true },
          ...(isAdmin
            ? [{
                title: '操作',
                render: (_: unknown, x: Patient) => (
                  <Button size="small" icon={<SwapOutlined />} onClick={() => openMerge(x.id)}>
                    合并到此档案
                  </Button>
                ),
              }]
            : []),
        ]}
      />

      <Drawer title="建立患者档案" width={520} open={open} onClose={() => setOpen(false)} destroyOnClose>
        <Form layout="vertical" onFinish={save} initialValues={{ gender: '男', age: 30, allergies: '无' }}>
          <Form.Item label="姓名" name="name" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item label="性别" name="gender"><Radio.Group options={genderOptions} /></Form.Item>
          <Form.Item label="年龄" name="age" rules={[{ required: true }]}><InputNumber min={0} max={150} style={{ width: '100%' }} /></Form.Item>
          <Form.Item label="身份证号" name="id_card" rules={[{ required: true, min: 10 }]}><Input /></Form.Item>
          <Form.Item label="联系方式" name="phone" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item label="过敏史" name="allergies"><Input.TextArea /></Form.Item>
          <Form.Item label="既往病史" name="medical_history"><Input.TextArea /></Form.Item>
          <Button htmlType="submit" type="primary">保存档案</Button>
        </Form>
      </Drawer>

      <Modal
        title="合并重复患者档案"
        open={mergeOpen}
        onCancel={() => setMergeOpen(false)}
        width={860}
        footer={[
          <Button key="cancel" onClick={() => setMergeOpen(false)}>取消</Button>,
          <Button key="preview" type="dashed" loading={previewLoading} onClick={() => void loadPreview()}>
            先显示病历数量
          </Button>,
          <Popconfirm
            key="confirm"
            title="确认执行档案合并？"
            description="系统会再次校验病历和处方数量；任一关联资料未转过去时将整体回滚。"
            onConfirm={() => void confirmMerge()}
            disabled={!preview || mergeLoading}
          >
            <Button type="primary" danger disabled={!preview} loading={mergeLoading}>
              确认合并
            </Button>
          </Popconfirm>,
        ]}
      >
        <Alert
          type="info"
          showIcon
          style={{ marginBottom: 16 }}
          message="第一套为保留档案；第二套合并后不会出现在患者检索和病历选择中，档案号、合并人、时间和原因会永久留档。"
        />
        <Form layout="vertical">
          <Form.Item label="保留档案（合并后继续使用）" required>
            <Select
              showSearch
              placeholder="搜索并选择保留档案"
              value={keepPatientId}
              onChange={(v) => { setKeepPatientId(v); setPreview(null); }}
              optionFilterProp="label"
              options={patientOptions}
            />
          </Form.Item>
          <Form.Item label="被合并档案（身份证录错的一套）" required>
            <Select
              showSearch
              placeholder="搜索并选择被合并档案"
              value={mergedPatientId}
              onChange={(v) => { setMergedPatientId(v); setPreview(null); }}
              optionFilterProp="label"
              options={patientOptions.filter((x) => x.value !== keepPatientId)}
            />
          </Form.Item>
          <Form.Item label="合并原因" required>
            <Input.TextArea
              rows={3}
              maxLength={1000}
              showCount
              value={reason}
              onChange={(e) => { setReason(e.target.value); setPreview(null); }}
              placeholder="例如：同一患者身份证号录入错误，重复建档。"
            />
          </Form.Item>
        </Form>

        {preview && (
          <Space direction="vertical" size={12} style={{ width: '100%' }}>
            <Card size="small" title="合并前数量确认">
              <Descriptions bordered size="small" column={3}>
                <Descriptions.Item label="保留档案病历">{preview.keep.record_count}</Descriptions.Item>
                <Descriptions.Item label="被合并档案病历">{preview.merged.record_count}</Descriptions.Item>
                <Descriptions.Item label="合并后病历总数">{preview.total_record_count}</Descriptions.Item>
                <Descriptions.Item label="保留档案处方">{preview.keep.prescription_count}</Descriptions.Item>
                <Descriptions.Item label="被合并档案处方">{preview.merged.prescription_count}</Descriptions.Item>
                <Descriptions.Item label="合并后处方总数">
                  {preview.keep.prescription_count + preview.merged.prescription_count}
                </Descriptions.Item>
              </Descriptions>
            </Card>
            <Card size="small" title="保留档案">
              <Descriptions size="small" column={2}>
                <Descriptions.Item label="档案号">{preview.keep.record_no}</Descriptions.Item>
                <Descriptions.Item label="姓名">{preview.keep.name}</Descriptions.Item>
                <Descriptions.Item label="身份证号">{preview.keep.id_card}</Descriptions.Item>
                <Descriptions.Item label="手机号">{preview.keep.phone}</Descriptions.Item>
              </Descriptions>
            </Card>
            <Card size="small" title="被合并档案（合并后隐藏）">
              <Descriptions size="small" column={2}>
                <Descriptions.Item label="档案号">{preview.merged.record_no}</Descriptions.Item>
                <Descriptions.Item label="姓名">{preview.merged.name}</Descriptions.Item>
                <Descriptions.Item label="身份证号">{preview.merged.id_card}</Descriptions.Item>
                <Descriptions.Item label="手机号">{preview.merged.phone}</Descriptions.Item>
              </Descriptions>
            </Card>
          </Space>
        )}
      </Modal>

      <Drawer title="患者档案合并记录" width={920} open={logsOpen} onClose={() => setLogsOpen(false)}>
        <Table
          rowKey="id"
          dataSource={logs}
          pagination={{ pageSize: 10 }}
          columns={[
            { title: '合并时间', dataIndex: 'merged_at', render: (v) => new Date(v).toLocaleString() },
            { title: '合并人', dataIndex: 'operator_name' },
            { title: '保留档案号', dataIndex: 'keep_record_no' },
            { title: '被合并档案号', dataIndex: 'merged_record_no' },
            { title: '病历数', dataIndex: 'moved_record_count' },
            { title: '处方数', dataIndex: 'moved_prescription_count' },
            { title: '原因', dataIndex: 'reason', width: 220 },
          ]}
        />
      </Drawer>
    </Card>
  );
}
