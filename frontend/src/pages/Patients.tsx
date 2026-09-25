import { HistoryOutlined, MergeCellsOutlined, PlusOutlined, SearchOutlined } from '@ant-design/icons';
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
  Radio,
  Select,
  Space,
  Table,
  Tag,
  Typography,
  message,
} from 'antd';
import { useEffect, useMemo, useState } from 'react';
import { api } from '../api/client';
import { useAuth } from '../store/auth';
import type { Patient, PatientMergePreview, PatientMergeRecord } from '../types';

const genderOptions = ['男', '女', '其他'];

const formatPatient = (p: Patient) =>
  `${p.name} · ${p.record_no} · ${p.id_card}`;

export default function Patients() {
  const { user } = useAuth();
  const [items, setItems] = useState<Patient[]>([]);
  const [createOpen, setCreateOpen] = useState(false);
  const [mergeOpen, setMergeOpen] = useState(false);
  const [historyOpen, setHistoryOpen] = useState(false);
  const [mergeHistory, setMergeHistory] = useState<PatientMergeRecord[]>([]);
  const [q, setQ] = useState('');
  const [retainedId, setRetainedId] = useState<number>();
  const [mergedId, setMergedId] = useState<number>();
  const [reason, setReason] = useState('');
  const [preview, setPreview] = useState<PatientMergePreview | null>(null);
  const [previewLoading, setPreviewLoading] = useState(false);
  const [submitLoading, setSubmitLoading] = useState(false);

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

  const patientOptions = useMemo(
    () => items.map((x) => ({ value: x.id, label: formatPatient(x) })),
    [items],
  );

  const save = async (v: object) => {
    try {
      await api.post('/patients', v);
      message.success('患者档案已建立');
      setCreateOpen(false);
      void load();
    } catch (e) {
      message.error((e as Error).message);
    }
  };

  const resetMerge = () => {
    setMergeOpen(false);
    setRetainedId(undefined);
    setMergedId(undefined);
    setReason('');
    setPreview(null);
  };

  const loadPreview = async () => {
    if (!retainedId || !mergedId) {
      message.warning('请先选择保留档案和被合并档案');
      return;
    }
    if (retainedId === mergedId) {
      message.warning('两套档案不能相同');
      return;
    }
    setPreviewLoading(true);
    try {
      const r = await api.post<PatientMergePreview>('/admin/patient-merges/preview', {
        retained_patient_id: retainedId,
        merged_patient_id: mergedId,
      });
      setPreview(r.data);
    } catch (e) {
      message.error((e as Error).message);
    } finally {
      setPreviewLoading(false);
    }
  };

  const loadHistory = async () => {
    try {
      const r = await api.get<PatientMergeRecord[]>('/admin/patient-merges');
      setMergeHistory(r.data);
      setHistoryOpen(true);
    } catch (e) {
      message.error((e as Error).message);
    }
  };

  const confirmMerge = () => {
    if (!preview || !reason.trim() || reason.trim().length < 4) {
      message.warning('请填写不少于 4 个字的合并原因');
      return;
    }
    Modal.confirm({
      title: '确认执行档案合并？',
      content: `将把 ${preview.merged.record_no} 的 ${preview.will_move_records} 份病历、${preview.will_move_prescriptions} 张处方转到 ${preview.retained.record_no}。被合并档案会从检索和选择中消失。`,
      okText: '确认合并',
      okType: 'danger',
      cancelText: '取消',
      onOk: async () => {
        setSubmitLoading(true);
        try {
          await api.post('/admin/patient-merges', {
            retained_patient_id: retainedId,
            merged_patient_id: mergedId,
            reason: reason.trim(),
          });
          message.success('档案已合并，关联病历和处方已转入保留档案');
          resetMerge();
          await load();
        } catch (e) {
          message.error((e as Error).message);
        } finally {
          setSubmitLoading(false);
        }
      },
    });
  };

  return (
    <Card
      className="page-card"
      title={<Typography.Title level={3} className="page-title">患者档案</Typography.Title>}
      extra={
        <Space wrap>
          <Input
            placeholder="姓名 / 身份证 / 手机号"
            value={q}
            onChange={(e) => setQ(e.target.value)}
            onPressEnter={load}
            suffix={<SearchOutlined />}
          />
          <Button onClick={load}>检索</Button>
          {user?.role === 'admin' && (
            <>
              <Button icon={<MergeCellsOutlined />} onClick={() => setMergeOpen(true)}>
                合并重复档案
              </Button>
              <Button icon={<HistoryOutlined />} onClick={loadHistory}>
                合并记录
              </Button>
            </>
          )}
          <Button type="primary" icon={<PlusOutlined />} onClick={() => setCreateOpen(true)}>
            新建档案
          </Button>
        </Space>
      }
    >
      <Table
        rowKey="id"
        dataSource={items}
        columns={[
          { title: '档案编号', dataIndex: 'record_no' },
          { title: '姓名', dataIndex: 'name' },
          {
            title: '性别',
            dataIndex: 'gender',
            render: (v) => <Tag color={v === '男' ? 'blue' : 'magenta'}>{v}</Tag>,
          },
          { title: '年龄', dataIndex: 'age' },
          { title: '身份证号', dataIndex: 'id_card' },
          { title: '联系方式', dataIndex: 'phone' },
          { title: '病历数', dataIndex: 'record_count', width: 90 },
          { title: '处方数', dataIndex: 'prescription_count', width: 90 },
          { title: '过敏史', dataIndex: 'allergies', ellipsis: true },
          { title: '既往病史', dataIndex: 'medical_history', ellipsis: true },
        ]}
      />

      <Drawer title="建立患者档案" width={520} open={createOpen} onClose={() => setCreateOpen(false)} destroyOnClose>
        <Form
          layout="vertical"
          onFinish={save}
          initialValues={{ gender: '男', age: 30, allergies: '无' }}
        >
          <Form.Item label="姓名" name="name" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item label="性别" name="gender"><Radio.Group options={genderOptions} /></Form.Item>
          <Form.Item label="年龄" name="age" rules={[{ required: true }]}>
            <InputNumber min={0} max={150} style={{ width: '100%' }} />
          </Form.Item>
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
        onCancel={resetMerge}
        width={720}
        footer={[
          <Button key="cancel" onClick={resetMerge}>取消</Button>,
          <Button key="preview" loading={previewLoading} onClick={loadPreview}>显示病历数量</Button>,
          <Button
            key="submit"
            type="primary"
            danger
            loading={submitLoading}
            disabled={!preview}
            onClick={confirmMerge}
          >
            确认合并
          </Button>,
        ]}
      >
        <Alert
          type="info"
          showIcon
          message="先检索并选择保留档案、被合并档案；确认后仅移动全部病历和处方。任一关联资料未成功移动时，本次合并会整体回滚。"
          style={{ marginBottom: 16 }}
        />
        <Form layout="vertical">
          <Form.Item label="保留档案（合并后继续使用）" required>
            <Select
              showSearch
              placeholder="检索后选择保留档案"
              value={retainedId}
              options={patientOptions.filter((x) => x.value !== mergedId)}
              filterOption={(input, option) => (option?.label ?? '').toLowerCase().includes(input.toLowerCase())}
              onChange={(v) => { setRetainedId(v); setPreview(null); }}
            />
          </Form.Item>
          <Form.Item label="被合并档案（合并后从检索和选择中消失）" required>
            <Select
              showSearch
              placeholder="检索后选择被合并档案"
              value={mergedId}
              options={patientOptions.filter((x) => x.value !== retainedId)}
              filterOption={(input, option) => (option?.label ?? '').toLowerCase().includes(input.toLowerCase())}
              onChange={(v) => { setMergedId(v); setPreview(null); }}
            />
          </Form.Item>
          <Form.Item label="合并原因" required>
            <Input.TextArea
              rows={3}
              value={reason}
              maxLength={500}
              showCount
              placeholder="如：身份证号录入错误，确认为同一患者"
              onChange={(e) => setReason(e.target.value)}
            />
          </Form.Item>
        </Form>

        {preview && (
          <Descriptions bordered size="small" column={2} title="合并前核对">
            <Descriptions.Item label="保留档案" span={2}>
              {preview.retained.name} · {preview.retained.record_no}
            </Descriptions.Item>
            <Descriptions.Item label="保留档案病历">
              {preview.retained.record_count}
            </Descriptions.Item>
            <Descriptions.Item label="保留档案处方">
              {preview.retained.prescription_count}
            </Descriptions.Item>
            <Descriptions.Item label="被合并档案" span={2}>
              {preview.merged.name} · {preview.merged.record_no}
            </Descriptions.Item>
            <Descriptions.Item label="将移动病历">
              {preview.will_move_records}
            </Descriptions.Item>
            <Descriptions.Item label="将移动处方">
              {preview.will_move_prescriptions}
            </Descriptions.Item>
          </Descriptions>
        )}
      </Modal>

      <Drawer title="档案合并记录" width={860} open={historyOpen} onClose={() => setHistoryOpen(false)}>
        <Table
          rowKey="id"
          dataSource={mergeHistory}
          pagination={false}
          columns={[
            {
              title: '合并时间',
              dataIndex: 'created_at',
              width: 170,
              render: (v) => new Date(v).toLocaleString(),
            },
            { title: '保留档案号', dataIndex: 'retained_record_no' },
            { title: '被合并档案号', dataIndex: 'merged_record_no' },
            { title: '合并人', dataIndex: 'merged_by_name' },
            { title: '病历', dataIndex: 'medical_record_count', width: 70 },
            { title: '处方', dataIndex: 'prescription_count', width: 70 },
            { title: '原因', dataIndex: 'reason', width: 220, ellipsis: true },
          ]}
        />
      </Drawer>
    </Card>
  );
}
